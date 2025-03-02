package dataloader

import (
	"context"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.EventDriver items
//
//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetEventDrivers(
	ctx context.Context,
	eventId int,
) ([]*model.EventDriver, []error) {
	thunk := i.driversByEventLoader.Load(ctx, storage.IntKey(eventId))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event driver data: %v", err)
		return nil, nil
	}
	//nolint:errcheck // we are sure that the type is correct
	ret := result.([]*model.EventDriver)
	return ret, nil
}

//nolint:whitespace // editor/linter issue
func (i *DataLoader) CollectDriversByEventEntry(
	ctx context.Context,
	eventId int,
) ([]*model.EventDriver, []error) {
	thunk := i.driversByEventEntryLoader.Load(ctx, storage.IntKey(eventId))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event driver data: %v", err)
		return nil, nil
	}
	//nolint:errcheck // we are sure that the type is correct
	ret := result.([]*model.EventDriver)
	return ret, nil
}

//nolint:whitespace // editor/linter issue
func (i *DataLoader) CollectDriversByTeam(
	ctx context.Context,
	teamId int,
) ([]*model.EventDriver, []error) {
	thunk := i.driversByTeamLoader.Load(ctx, storage.IntKey(teamId))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event driver data: %v", err)
		return nil, nil
	}
	//nolint:errcheck // we are sure that the type is correct
	ret := result.([]*model.EventDriver)
	return ret, nil
}
