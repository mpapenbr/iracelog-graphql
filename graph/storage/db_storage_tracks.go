package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/tracks"
)

// contains implementations of storage interface that return a model.Track items
//
//nolint:whitespace // editor/linter issue
func (db *DBStorage) GetAllTracks(
	ctx context.Context,
	limit *int,
	offset *int,
	sort []*model.TrackSortArg,
) ([]*model.Track, error) {
	var result []*model.Track

	dbTrackSortArg := convertTrackSortArgs(sort)
	tracks, err := tracks.GetAll(
		db.executor,
		internal.DBPageable{Limit: limit, Offset: offset, Sort: dbTrackSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, track := range tracks {
			result = append(result, convertDBTrackToModel(track))
		}
	}
	return result, err
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) GetTracksByKeys(
	ctx context.Context,
	ids dataloader.Keys,
) map[string]*model.Track {
	intIDs := IntKeysToSlice(ids)
	result := map[string]*model.Track{}

	tracks, _ := tracks.GetByIDs(db.executor, intIDs)

	// convert the internal database Track to the GraphQL-Track
	for _, track := range tracks {
		result[IntKey(track.ID).String()] = convertDBTrackToModel(track)
	}
	return result
}
