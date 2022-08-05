package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/graph-gophers/dataloader"
	gopher_dataloader "github.com/graph-gophers/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
)

type ctxKey string

const loadersKey = ctxKey("dataloader")

type DataLoader struct {
	trackLoader *dataloader.Loader
	eventLoader *dataloader.Loader
}

// GetUser wraps the User dataloader for efficient retrieval by user ID
func (i *DataLoader) GetTrack(ctx context.Context, trackID int) (*model.Track, error) {
	thunk := i.trackLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", trackID)))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.Track), nil
}

func Middleware(loaders *DataLoader, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *DataLoader {
	return ctx.Value(loadersKey).(*DataLoader)
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(db storage.Storage) *DataLoader {
	// define the data loader
	tracks := &trackBatcher{db: db}
	loaders := &DataLoader{
		trackLoader: dataloader.NewBatchedLoader(tracks.get),
	}
	return loaders
}

// trackBatcher wraps storage and provides a "get" method for the track dataloader
type trackBatcher struct {
	db storage.Storage
}

// get implements the dataloader for finding many tracks by Id and returns
// them in the order requested
func (t *trackBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	fmt.Printf("dataloader.trackBatcher.get, tracks: [%s]\n", strings.Join(keys.Keys(), ","))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[int]int, len(keys))
	// collect the keys to search for
	var trackIDs []int
	for ix, key := range keys {
		id, _ := strconv.Atoi(key.String())
		trackIDs = append(trackIDs, id)
		keyOrder[id] = ix
	}
	// search for those users

	dbRecords, err := t.db.GetTracks(ctx, trackIDs)
	// if DB error, return
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for _, record := range dbRecords {
		ix, ok := keyOrder[record.ID]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, record.ID)
		}
	}
	// fill array positions with errors where not found in DB
	for userID, ix := range keyOrder {
		err := fmt.Errorf("track not found %d", userID)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}
