package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"github.com/mpapenbr/iracelog-graphql/graph"
	"github.com/mpapenbr/iracelog-graphql/graph/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/generated"
	"github.com/mpapenbr/iracelog-graphql/graph/resolver"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := storage.NewDbStorage()
	// _ := database.InitDB()

	graphResolver := resolver.NewResolver(db)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))
	dataloaderSrv := dataloader.Middleware(db, srv)

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "*"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	// CORS handling
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", dataloaderSrv)

	log.Printf("iRacelog GraphQL service %s", graph.Version)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
