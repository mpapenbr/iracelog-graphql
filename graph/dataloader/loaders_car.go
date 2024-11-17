package dataloader

import (
	"context"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.Car items

func (i *DataLoader) GetEventCars(
	ctx context.Context,
	eventId int,
) ([]*model.Car, []error) {
	thunk := i.eventCarsLoader.Load(ctx, storage.IntKey(eventId))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil, nil
	}
	ret := result.([]*model.Car)
	return ret, nil
}

func (i *DataLoader) GetEventEntryCar(
	ctx context.Context,
	eventEntryId int,
) (*model.Car, []error) {
	thunk := i.eventEntryCarLoader.Load(ctx, storage.IntKey(eventEntryId))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event car data: %v", err)
		return nil, nil
	}
	ret := result.(*model.Car)
	return ret, nil
}
