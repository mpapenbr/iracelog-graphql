package dataloader

import (
	"context"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.EventEntry items

func (i *DataLoader) GetEventEntriesForIds(
	ctx context.Context,
	ids []int,
) ([]*model.EventEntry, []error) {
	eventKeys := storage.NewKeysFromInts(ids)
	thunkMany := i.eventEntriesByIdsLoader.LoadMany(ctx, eventKeys)
	result, err := thunkMany()
	if err != nil {
		return nil, err
	}

	// hmm, this copy bothers me, but my "wish-statement" return result.([]*model.Event) doesn't work
	ret := make([]*model.EventEntry, len(result))
	for i, v := range result {
		ret[i] = v.(*model.EventEntry)
	}
	return ret, err
}

func (i *DataLoader) GetEventEntries(
	ctx context.Context,
	eventId int,
) ([]*model.EventEntry, []error) {
	thunk := i.eventEntriesByEventLoader.Load(ctx, storage.IntKey(eventId))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event entry data: %v", err)
		return nil, nil
	}
	ret := result.([]*model.EventEntry)
	return ret, nil
}
