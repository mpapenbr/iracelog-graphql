package dataloader

import (
	"context"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.EventEntry items
//
//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetEventEntriesForIDs(
	ctx context.Context,
	ids []int,
) ([]*model.EventEntry, []error) {
	eventKeys := storage.NewKeysFromInts(ids)
	thunkMany := i.eventEntriesByIDsLoader.LoadMany(ctx, eventKeys)
	result, err := thunkMany()
	if err != nil {
		return nil, err
	}

	// hmm, this copy bothers me, but my "wish-statement"

	ret := make([]*model.EventEntry, len(result))
	for i, v := range result {
		//nolint:errcheck // we are sure that the type is correct
		ret[i] = v.(*model.EventEntry)
	}
	return ret, err
}

//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetEventEntries(
	ctx context.Context,
	eventID int,
) ([]*model.EventEntry, []error) {
	thunk := i.eventEntriesByEventLoader.Load(ctx, storage.IntKey(eventID))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event entry data: %v", err)
		return nil, nil
	}
	//nolint:errcheck // we are sure that the type is correct
	ret := result.([]*model.EventEntry)
	return ret, nil
}
