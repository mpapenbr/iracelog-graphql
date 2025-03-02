package dataloader

import (
	"context"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.Track items

// GetTrack wraps the Track dataloader for efficient retrieval by track ID
//
//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetTrack(
	ctx context.Context,
	trackID int,
) (*model.Track, error) {
	thunk := i.trackLoader.Load(ctx, storage.IntKey(trackID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // we are sure that the type is correct
	return result.(*model.Track), nil
}

//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetTracks(
	ctx context.Context,
	trackIds []int,
) ([]*model.Track, []error) {
	trackKeys := storage.NewKeysFromInts(trackIds)
	thunkMany := i.trackLoader.LoadMany(ctx, trackKeys)
	result, err := thunkMany()
	if err != nil {
		return nil, err
	}

	// hmm, this copy bothers me, but my "wish-statement"
	// return result.([]*model.Track) doesn't work
	ret := make([]*model.Track, len(result))
	for i, v := range result {
		//nolint:errcheck // we are sure that the type is correct
		ret[i] = v.(*model.Track)
	}
	return ret, err
}
