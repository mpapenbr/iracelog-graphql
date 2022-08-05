package storage

import (
	"context"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
)

type Storage interface {

	// GetTracks accepts many user IDs and returns an array of matching Tracks
	GetTracks(ctx context.Context, ids []int) ([]*model.Track, error)
	// GetEvents accepts many event IDs and returns an array of matching Tracks
	GetEvents(ctx context.Context, ids []int) ([]*model.Event, error)

	// GetAllTracks lists all Tracks in the database
	GetAllTracks(ctx context.Context) ([]*model.Track, error)
	// GetAllEvents lists all Events in the database
	GetAllEvents(ctx context.Context) ([]*model.Event, error)
}
