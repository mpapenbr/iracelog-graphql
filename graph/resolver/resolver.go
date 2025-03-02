package resolver

//go:generate go run github.com/99designs/gqlgen generate

import "github.com/mpapenbr/iracelog-graphql/graph/storage"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app,
// add any dependencies you require here.

type Resolver struct {
	// db is an interface for reading/writing to the datastore
	db storage.Storage
}

// NewResolver returns a Resolver
func NewResolver(db storage.Storage) *Resolver {
	output := &Resolver{db: db}
	return output
}
