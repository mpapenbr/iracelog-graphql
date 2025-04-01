package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/car/team"
)

// contains implementations of storage interface that return a model.EventTeam items
//
//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectTeamByEventEntry(
	ctx context.Context,
	eventEntryIDs dataloader.Keys,
) map[string]*model.EventTeam {
	res, _ := team.GetTeamsByEventEntry(db.executor, IntKeysToSlice(eventEntryIDs))
	ret := map[string]*model.EventTeam{}
	for k, d := range res {
		key := IntKey(k).String()

		ed := convertDBCarTeamToModel(d)

		ret[key] = ed
	}
	return ret
}
