//nolint:dupl // false positive
package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/car/entry"
)

// contains implementations of storage interface that return a model.EventEntry items
//
//nolint:whitespace // editor/linter issue
func (db *DbStorage) CollectEventEntries(
	ctx context.Context,
	eventIds dataloader.Keys,
) map[string][]*model.EventEntry {
	res, _ := entry.GetEventEntriesByEventId(db.pool, IntKeysToSlice(eventIds))
	ret := map[string][]*model.EventEntry{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventEntry, len(v))
		for i, d := range v {
			ed[i] = convertDbEventEntryToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

//nolint:whitespace // editor/linter issue
func (db *DbStorage) CollectEventEntriesById(
	ctx context.Context,
	ids dataloader.Keys,
) map[string]*model.EventEntry {
	res, _ := entry.GetEventEntriesByIds(db.pool, IntKeysToSlice(ids))
	ret := map[string]*model.EventEntry{}
	for k, d := range res {
		key := IntKey(k).String()
		ed := convertDbEventEntryToModel(d)
		ret[key] = ed
	}
	return ret
}
