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
func (db *DbStorage) CollectEventDrivers(
	ctx context.Context,
	eventIds dataloader.Keys,
) map[string][]*model.EventDriver {
	res, _ := driver.GetEventDrivers(ctx, db.executor, IntKeysToSlice(eventIds))
	ret := map[string][]*model.EventDriver{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventDriver, len(v))
		for i, d := range v {
			ed[i] = convertDbCarDriverToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

//nolint:dupl,whitespace // false positive
func (db *DbStorage) CollectDriversByEventEntry(
	ctx context.Context,
	eventEntryIds dataloader.Keys,
) map[string][]*model.EventDriver {
	res, _ := driver.GetDriversByEventEntry(
		ctx,
		db.executor,
		IntKeysToSlice(eventEntryIds))
	ret := map[string][]*model.EventDriver{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventDriver, len(v))
		for i, d := range v {
			ed[i] = convertDbCarDriverToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

//nolint:dupl,whitespace // false positive
func (db *DbStorage) CollectDriversByTeam(
	ctx context.Context,
	teamIds dataloader.Keys,
) map[string][]*model.EventDriver {
	res, _ := driver.GetDriversByTeam(ctx, db.executor, IntKeysToSlice(teamIds))
	ret := map[string][]*model.EventDriver{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.EventDriver, len(v))
		for i, d := range v {
			ed[i] = convertDbCarDriverToModel(d)
		}
		ret[key] = ed
	}
	return ret
}
