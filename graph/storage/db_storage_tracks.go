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
func (db *DbStorage) GetAllTracks(
	ctx context.Context,
	limit *int,
	offset *int,
	sort []*model.TrackSortArg,
) ([]*model.Track, error) {
	var result []*model.Track

	dbTrackSortArg := convertTrackSortArgs(sort)
	tracks, err := tracks.GetALl(
		db.pool,
		internal.DbPageable{Limit: limit, Offset: offset, Sort: dbTrackSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, track := range tracks {
			result = append(result, convertDbTrackToModel(track))
		}
	}
	return result, err
}

//nolint:whitespace // editor/linter issue
func (db *DbStorage) GetTracksByKeys(
	ctx context.Context,
	ids dataloader.Keys,
) map[string]*model.Track {
	intIds := IntKeysToSlice(ids)
	result := map[string]*model.Track{}

	tracks, _ := tracks.GetByIds(db.pool, intIds)

	// convert the internal database Track to the GraphQL-Track
	for _, track := range tracks {
		result[IntKey(track.ID).String()] = convertDbTrackToModel(track)
	}
	return result
}
