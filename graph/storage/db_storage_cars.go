//nolint:dupl // false positive
package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/car/car"
)

// contains implementations of storage interface that return a model.Car items
//
//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectEventCars(
	ctx context.Context,
	eventIDs dataloader.Keys,
) map[string][]*model.Car {
	res, _ := car.GetEventCars(db.executor, IntKeysToSlice(eventIDs))
	ret := map[string][]*model.Car{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.Car, len(v))
		for i, d := range v {
			ed[i] = convertDBCarToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectCarsByEventEntry(
	ctx context.Context,
	eventEntryIDs dataloader.Keys,
) map[string]*model.Car {
	res, _ := car.GetEventEntryCars(db.executor, IntKeysToSlice(eventEntryIDs))
	ret := map[string]*model.Car{}
	for k, d := range res {
		key := IntKey(k).String()
		ed := convertDBCarToModel(d)
		ret[key] = ed
	}
	return ret
}
