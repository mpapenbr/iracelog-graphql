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
func (db *DBStorage) CollectEventEntries(
	ctx context.Context,
	eventIDs dataloader.Keys,
) map[string][]*model.EventEntry {
	res, _ := entry.GetEventEntriesByEventID(db.executor, IntKeysToSlice(eventIDs))
	ret := map[string][]*model.EventEntry{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventEntry, len(v))
		for i, d := range v {
			ed[i] = convertDBEventEntryToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectEventEntriesByID(
	ctx context.Context,
	ids dataloader.Keys,
) map[string]*model.EventEntry {
	res, _ := entry.GetEventEntriesByIDs(db.executor, IntKeysToSlice(ids))
	ret := map[string]*model.EventEntry{}
	for k, d := range res {
		key := IntKey(k).String()
		ed := convertDBEventEntryToModel(d)
		ret[key] = ed
	}
	return ret
}
