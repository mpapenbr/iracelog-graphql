package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgx-contrib/pgxtrace"
	"github.com/spf13/cobra"

	"github.com/mpapenbr/iracelog-graphql/config"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
	"github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
	"github.com/mpapenbr/iracelog-graphql/internal/server"
	"github.com/mpapenbr/iracelog-graphql/internal/utils"
	"github.com/mpapenbr/iracelog-graphql/log"
)

var (
	supportTenants  bool
	defaultTenantID int = 1 // default tenant id (used if tenants are disabled)
)

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the GraphQL server",
		Run: func(cmd *cobra.Command, args []string) {
			// start the server
			startServer(cmd.Context())
		},
	}
	cmd.Flags().BoolVar(&supportTenants,
		"enable-tenants",
		false,
		"enables tenant support")
	cmd.Flags().IntVar(&defaultTenantID,
		"default-tenant-id",
		defaultTenantID,
		"id of the internal default tenant (will be used if tenants are disabled)")
	return cmd
}

type graphqlServer struct {
	ctx            context.Context
	log            *log.Logger
	pool           *pgxpool.Pool
	supportTenants bool
}

func startServer(ctx context.Context) {
	srv := &graphqlServer{
		ctx:            ctx,
		supportTenants: supportTenants,
	}
	srv.SetupLogger()
	srv.waitForRequiredServices()
	srv.SetupDBPgx()

	if err := srv.Start(); err != nil {
		srv.log.Error("error starting server", log.ErrorField(err))
	}
}

func (s *graphqlServer) SetupLogger() {
	s.log = log.GetFromContext(s.ctx).Named("server")
}

func (s *graphqlServer) SetupDBPgx() {
	pgTracer := pgxtrace.CompositeQueryTracer{
		postgres.NewMyTracer(log.GetFromContext(s.ctx).Named("sql"), log.DebugLevel),
	}
	//nolint:gocritic // will be used later
	// if config.EnableTelemetry {
	// 	var err error
	// 	if s.telemetry, err = config.SetupTelemetry(context.Background()); err == nil {
	// 		pgTracer = append(pgTracer, postgres.NewOtlpTracer())
	// 	} else {
	// 		s.log.Warn("Could not setup db telemetry", log.ErrorField(err))
	// 	}
	// }

	pgOptions := []postgres.PoolConfigOption{
		postgres.WithTracer(pgTracer),
	}
	s.pool = postgres.InitWithURL(
		config.DB,
		pgOptions...,
	)
	s.log.Info("PgxPool initialized")
}

func (s *graphqlServer) Start() error {
	ch := make(chan error, 2)

	s.log.Info("Starting server")
	go func() {
		storageOpts := []storage.DBStorageOption{}

		myStorage := storage.NewDBStorage(s.pool, storageOpts...)
		opts := []server.Option{
			server.WithContext(s.ctx),
			server.WithLogger(log.GetFromContext(s.ctx).Named("gql")),
			server.WithStorage(myStorage),
			server.WithTenantResolver(func(r *http.Request) (int, error) {
				return defaultTenantID, nil
			}),
		}
		if config.Addr != "" {
			opts = append(opts, server.WithAddr(config.Addr))
		}
		if s.supportTenants {
			s.log.Info("enabled tenant support")
			opts = append(opts, server.WithRequestBasedTenantResolver())
		}

		gqlServer := server.NewServer(opts...)
		ch <- gqlServer.Start()
	}()

	s.log.Debug("Wait for signal or server termination")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	select {
	case v := <-sigChan:
		s.log.Debug("Got signal", log.Any("signal", v))
	case err := <-ch:
		s.log.Debug("Server terminated", log.ErrorField(err))
	}
	s.pool.Close()
	s.log.Info("Server stopped")
	return nil
}

func (s *graphqlServer) waitForRequiredServices() {
	var err error
	wg := sync.WaitGroup{}
	checkTCP := func(addr string) {
		if err = utils.WaitForTCP(addr, config.WaitForServices); err != nil {
			s.log.Fatal("required services not ready", log.ErrorField(err))
		}
		wg.Done()
	}

	if postgresAddr := utils.ExtractFromDBURL(config.DB); postgresAddr != "" {
		wg.Add(1)
		go checkTCP(postgresAddr)
	}

	s.log.Debug("Waiting for connection checks to return")
	wg.Wait()
	s.log.Debug("Required services are available")
}
