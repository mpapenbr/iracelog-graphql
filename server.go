package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", dataloaderSrv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
