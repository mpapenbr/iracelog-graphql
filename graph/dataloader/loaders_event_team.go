package dataloader

import (
	"context"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

// contains implementations of DataLoader struct that return a model.EventTeam items

//nolint:whitespace // editor/linter issue
func (i *DataLoader) GetTeamByEventEntry(
	ctx context.Context,
	eventEntryID int,
) (*model.EventTeam, []error) {
	thunk := i.teamByEventEntryLoader.Load(ctx, storage.IntKey(eventEntryID))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading event team data: %v", err)
		return nil, nil
	}
	//nolint:errcheck // we are sure that the type is correct
	ret := result.(*model.EventTeam)
	return ret, nil
}
