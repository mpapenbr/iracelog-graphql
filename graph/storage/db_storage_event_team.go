package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/car/team"
)

// contains implementations of storage interface that return a model.EventTeam items

func (db *DbStorage) CollectTeamByEventEntry(
	ctx context.Context,
	eventEntryIds dataloader.Keys,
) map[string]*model.EventTeam {
	res, _ := team.GetTeamsByEventEntry(db.pool, IntKeysToSlice(eventEntryIds))
	ret := map[string]*model.EventTeam{}
	for k, d := range res {
		key := IntKey(k).String()

		ed := convertDbCarTeamToModel(d)

		ret[key] = ed
	}
	return ret
}
