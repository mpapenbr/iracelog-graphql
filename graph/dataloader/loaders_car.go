package dataloader

import (
	"context"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.Car items
//
//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetEventCars(
	ctx context.Context,
	eventID int,
) ([]*model.Car, []error) {
	thunk := i.eventCarsLoader.Load(ctx, storage.IntKey(eventID))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil, nil
	}
	//nolint:errcheck // we are sure that the type is correct
	ret := result.([]*model.Car)
	return ret, nil
}

//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetEventEntryCar(
	ctx context.Context,
	eventEntryID int,
) (*model.Car, []error) {
	thunk := i.eventEntryCarLoader.Load(ctx, storage.IntKey(eventEntryID))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event car data: %v", err)
		return nil, nil
	}
	//nolint:errcheck // we are sure that the type is correct
	ret := result.(*model.Car)
	return ret, nil
}
