package tcpostgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// create a pg connection pool for the iracelog testdatabase
func SetupTestDb() *pgxpool.Pool {
	ctx := context.Background()
	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		log.Fatal(err)
	}
	container, err := SetupPostgres(ctx,
		WithPort(port.Port()),
		WithInitialDatabase("postgres", "password", "postgres"),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(1).WithStartupTimeout(5*time.Second)),
		WithName("iracelog-graphql-test"),
	)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n", container)
	containerPort, _ := container.MappedPort(ctx, port)
	host, _ := container.Host(ctx)
	dbUrl := fmt.Sprintf("postgresql://postgres:password@%s:%s/postgres", host, containerPort.Port())
	pool := database.InitWithUrl(dbUrl)
	return pool
}
