package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
	"github.com/mpapenbr/iracelog-graphql/internal/tracks"
)

type DbStorage struct {
	Storage
	pool *pgxpool.Pool
}

func NewDbStorage() *DbStorage {
	return &DbStorage{pool: database.InitDB()}
}

// tracks
func (db *DbStorage) GetAllTracks(ctx context.Context) ([]*model.Track, error) {

	var result []*model.Track

	tracks, err := tracks.GetALl(db.pool)
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, track := range tracks {
			result = append(result, &model.Track{ID: track.ID, Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length})
		}
	}
	return result, err
}
func (db *DbStorage) GetTracks(ctx context.Context, ids []int) ([]*model.Track, error) {

	var result []*model.Track

	tracks, err := tracks.GetByIds(db.pool, ids)
	log.Printf("Tracks: %v\n", tracks)
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, track := range tracks {
			result = append(result, &model.Track{ID: track.ID, Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length})
		}
	}
	return result, err
}

// events
func (db *DbStorage) GetAllEvents(ctx context.Context) ([]*model.Event, error) {

	var result []*model.Event

	events, err := events.GetALl(db.pool)
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, event := range events {
			result = append(result, &model.Event{ID: event.ID, Name: event.Name, Key: event.Key, TrackId: event.Info.TrackId})
		}
	}
	return result, err
}
