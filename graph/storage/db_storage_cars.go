package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/car/car"
)

// contains implementations of storage interface that return a model.Car items

func (db *DbStorage) CollectEventCars(
	ctx context.Context,
	eventIds dataloader.Keys,
) map[string][]*model.Car {
	res, _ := car.GetEventCars(db.pool, IntKeysToSlice(eventIds))
	ret := map[string][]*model.Car{}
	for k, v := range res {
		key := IntKey(k).String()
		ed := make([]*model.Car, len(v))
		for i, d := range v {
			ed[i] = convertDbCarToModel(d)
		}
		ret[key] = ed
	}
	return ret
}

func (db *DbStorage) CollectCarsByEventEntry(
	ctx context.Context,
	eventEntryIds dataloader.Keys,
) map[string]*model.Car {
	res, _ := car.GetEventEntryCars(db.pool, IntKeysToSlice(eventEntryIds))
	ret := map[string]*model.Car{}
	for k, d := range res {
		key := IntKey(k).String()
		ed := convertDbCarToModel(d)
		ret[key] = ed
	}
	return ret
}
