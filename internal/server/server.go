package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/mpapenbr/iracelog-graphql/graph/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/generated"
	"github.com/mpapenbr/iracelog-graphql/graph/resolver"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
	"github.com/mpapenbr/iracelog-graphql/log"
	"github.com/mpapenbr/iracelog-graphql/version"
)

type (
	Option func(*Server)

	Server struct {
		ctx  context.Context
		log  *log.Logger
		db   storage.Storage
		addr string
	}
)

func NewServer(opts ...Option) *Server {
	ret := &Server{
		addr: "localhost:8080",
		ctx:  context.Background(),
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func WithStorage(db storage.Storage) Option {
	return func(s *Server) {
		s.db = db
	}
}

func WithContext(arg context.Context) Option {
	return func(s *Server) {
		s.ctx = arg
	}
}

func WithLogger(arg *log.Logger) Option {
	return func(s *Server) {
		s.log = arg
	}
}

func WithAddr(arg string) Option {
	return func(s *Server) {
		s.addr = arg
	}
}

func (s *Server) Start() error {
	graphResolver := resolver.NewResolver(s.db)
	srv := handler.New(
		generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}),
	)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	dataloaderSrv := dataloader.Middleware(s.db, srv)

	router := chi.NewRouter()

	// add logger to context
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCtx := log.AddToContext(r.Context(), s.log)
			r = r.WithContext(newCtx)
			h.ServeHTTP(w, r)
		})
	})

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.log.Debug("Request", log.String("url", r.URL.String()))
			h.ServeHTTP(w, r)
		})
	})

	// CORS handling
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", dataloaderSrv)
	router.Handle("/healthz", healthzHandler())
	s.log.Info("iRacelog GraphQL service", log.String("version", version.FullVersion))
	s.log.Info("Listen", log.String("addr", s.addr))
	//nolint:gosec // tbd
	return http.ListenAndServe(s.addr, router)
}

func healthzHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Status string `json:"status"`
		}{Status: "ok"}

		w.Header().Set("Content-Type", "application/json")
		respData, _ := json.Marshal(data)
		if _, err := w.Write(respData); err != nil {
			log.Error("error writing healthz response", log.ErrorField(err))
		}
	})
}
