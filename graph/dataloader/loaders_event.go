package dataloader

import (
	"context"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.Event items
func (i *DataLoader) GetEvents(ctx context.Context, eventIDs []int) ([]*model.Event, []error) {
	eventKeys := storage.NewKeysFromInts(eventIDs)
	thunkMany := i.eventLoader.LoadMany(ctx, eventKeys)
	result, err := thunkMany()
	if err != nil {
		return nil, err
	}

	// hmm, this copy bothers me, but my "wish-statement" return result.([]*model.Event) doesn't work
	ret := make([]*model.Event, len(result))
	for i, v := range result {
		ret[i] = v.(*model.Event)
	}
	return ret, err
}

// GetEvent wraps the Event dataloader for efficient retrieval by event ID
func (i *DataLoader) GetEvent(ctx context.Context, eventID int) (*model.Event, error) {
	thunk := i.eventLoader.Load(ctx, storage.IntKey(eventID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.Event), nil
}

func (i *DataLoader) GetEventsForTrack(ctx context.Context, trackId int) []*model.Event {
	// thunk := i.eventsByTrackLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", trackId)))
	thunk := i.eventsByTrackLoader.Load(ctx, storage.IntKey(trackId))
	result, err := thunk()
	if err != nil {
		return nil
	}
	return result.([]*model.Event)
}
