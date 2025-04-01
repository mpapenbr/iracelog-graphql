package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/car/driver"
)

// contains implementations of storage interface that return a model.EventDriver items
//
//nolint:dupl,whitespace // false positive
func (db *DBStorage) CollectEventDrivers(
	ctx context.Context,
	eventIDs dataloader.Keys,
) map[string][]*model.EventDriver {
	res, _ := driver.GetEventDrivers(ctx, db.executor, IntKeysToSlice(eventIDs))
	ret := map[string][]*model.EventDriver{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventDriver, len(v))
		for i, d := range v {
			ed[i] = convertDBCarDriverToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

//nolint:dupl,whitespace // false positive
func (db *DBStorage) CollectDriversByEventEntry(
	ctx context.Context,
	eventEntryIDs dataloader.Keys,
) map[string][]*model.EventDriver {
	res, _ := driver.GetDriversByEventEntry(
		ctx,
		db.executor,
		IntKeysToSlice(eventEntryIDs))
	ret := map[string][]*model.EventDriver{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventDriver, len(v))
		for i, d := range v {
			ed[i] = convertDBCarDriverToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

//nolint:dupl,whitespace // false positive
func (db *DBStorage) CollectDriversByTeam(
	ctx context.Context,
	teamIDs dataloader.Keys,
) map[string][]*model.EventDriver {
	res, _ := driver.GetDriversByTeam(ctx, db.executor, IntKeysToSlice(teamIDs))
	ret := map[string][]*model.EventDriver{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventDriver, len(v))
		for i, d := range v {
			ed[i] = convertDBCarDriverToModel(d)
		}
		ret[key] = ed
	}
	return ret
}
