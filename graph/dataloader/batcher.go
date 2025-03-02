package dataloader

import (
	"context"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/log"
)

// this is the overall contract for this batcher:
// the collector has to provide a map which keys are built by Key.String()

type MapCollector[V any] func(ctx context.Context, keys dataloader.Keys) map[string]V

type genericMapBatcher[V any] struct {
	name      string
	collector MapCollector[V]
}

// get implements the dataloader for finding data by keys and returns
// them in the order requested keys provide a String() func
// that provides a unique value over all keys used in that batch
//
//nolint:whitespace // editor/linter issue
func (t *genericMapBatcher[V]) get(
	ctx context.Context,
	keys dataloader.Keys,
) []*dataloader.Result {
	l := log.GetFromContext(ctx)
	l.Debug(t.name, log.String("keys", strings.Join(keys.Keys(), ",")))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// collect the keys to search for

	for ix, key := range keys {
		id := key.String()
		keyOrder[id] = ix

	}
	// search for those keys

	dbRecords := t.collector(ctx, keys)

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for key, record := range dbRecords {
		ix, ok := keyOrder[key]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, key)
		}
	}
	// fill array positions with errors where not found in DB
	for key, ix := range keyOrder {
		err := fmt.Errorf("key not found %s", key)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	return results
}
