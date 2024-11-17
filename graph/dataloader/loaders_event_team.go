package dataloader

import (
	"context"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.EventTeam items

func (i *DataLoader) GetTeamByEventEntry(
	ctx context.Context,
	eventEntryId int,
) (*model.EventTeam, []error) {
	thunk := i.teamByEventEntryLoader.Load(ctx, storage.IntKey(eventEntryId))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event team data: %v", err)
		return nil, nil
	}
	ret := result.(*model.EventTeam)
	return ret, nil
}
