package dataloader

import (
	"context"
	"fmt"
	"log"
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
	trackLoader  *dataloader.Loader
	eventLoader  *dataloader.Loader
	driverLoader *dataloader.Loader
	teamLoader   *dataloader.Loader
}

// GetTrack wraps the Track dataloader for efficient retrieval by track ID
func (i *DataLoader) GetTrack(ctx context.Context, trackID int) (*model.Track, error) {
	thunk := i.trackLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", trackID)))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.Track), nil
}

func (i *DataLoader) GetTracks(ctx context.Context, trackIds []int) ([]*model.Track, []error) {
	trackKeys := make([]gopher_dataloader.Key, len(trackIds))
	for idx, val := range trackIds {
		trackKeys[idx] = gopher_dataloader.StringKey(fmt.Sprintf("%d", val))
	}
	thunkMany := i.trackLoader.LoadMany(ctx, trackKeys)
	result, err := thunkMany()
	if err != nil {
		return nil, err
	}

	// hmm, this copy bothers me, but my "wish-statement" return result.([]*model.Track) doesn't work
	ret := make([]*model.Track, len(result))
	for i, v := range result {
		ret[i] = v.(*model.Track)
	}
	return ret, err
}

func (i *DataLoader) GetEvents(ctx context.Context, eventIDs []int) ([]*model.Event, []error) {
	eventKeys := make([]gopher_dataloader.Key, len(eventIDs))
	for idx, val := range eventIDs {
		eventKeys[idx] = gopher_dataloader.StringKey(fmt.Sprintf("%d", val))
	}
	thunkMany := i.eventLoader.LoadMany(ctx, eventKeys)
	result, err := thunkMany()
	if err != nil {
		return nil, err
	}

	// hmm, this copy bothers me, but my "wish-statement" return result.([]*model.Event) doesn't work
	ret := make([]*model.Event, len(result))
	for i, v := range result {
		ret[i] = v.(*model.Event)
	}
	return ret, err
}

// GetEvent wraps the Event dataloader for efficient retrieval by event ID
func (i *DataLoader) GetEvent(ctx context.Context, eventID int) (*model.Event, error) {
	thunk := i.eventLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", eventID)))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.Event), nil
}

func (i *DataLoader) GetTeamDrivers(ctx context.Context, team string) ([]*model.Driver, []error) {
	thunk := i.driverLoader.Load(ctx, gopher_dataloader.StringKey(team))
	result, err := thunk()
	if err != nil {
		return nil, nil
	}
	return result.([]*model.Driver), nil
}

func (i *DataLoader) GetDriversTeams(ctx context.Context, driver string) ([]*model.Team, []error) {
	thunk := i.teamLoader.Load(ctx, gopher_dataloader.StringKey(driver))
	result, err := thunk()
	if err != nil {
		return nil, nil
	}
	return result.([]*model.Team), nil
}

func Middleware(db *storage.DbStorage, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := NewDataLoader(db) // we want fresh loaders on each request
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
		// loaders.driverLoader.ClearAll()
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
	events := &eventBatcher{db: db}
	drivers := &driverBatcher{db: db}
	teams := &teamBatcher{db: db}
	loaders := &DataLoader{
		trackLoader:  dataloader.NewBatchedLoader(tracks.get),
		eventLoader:  dataloader.NewBatchedLoader(events.get),
		driverLoader: dataloader.NewBatchedLoader(drivers.get),
		teamLoader:   dataloader.NewBatchedLoader(teams.get),
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
	log.Printf("dataloader.trackBatcher.get, tracks: [%s]\n", strings.Join(keys.Keys(), ","))
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

// events

// eventBatcher wraps storage and provides a "get" method for the track dataloader
type eventBatcher struct {
	db storage.Storage
}

// get implements the dataloader for finding many tracks by Id and returns
// them in the order requested
func (t *eventBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	for i, v := range keys {
		log.Printf("i:%v v:%v\n", i, v)
	}
	log.Printf("dataloader.eventBatcher.get, events: [%s]\n", strings.Join(keys.Keys(), ","))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[int]int, len(keys))
	// collect the keys to search for
	var eventIDs []int
	for ix, key := range keys {
		id, _ := strconv.Atoi(key.String())
		eventIDs = append(eventIDs, id)
		keyOrder[id] = ix
	}
	// search for those users

	dbRecords, err := t.db.GetEvents(ctx, eventIDs)
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
		err := fmt.Errorf("event not found %d", userID)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}

// drivers

// driverBatcher wraps storage and provides a "get" method for the driver dataloader
type driverBatcher struct {
	db storage.Storage
}

// get implements the dataloader for finding many tracks by Id and returns
// them in the order requested
func (t *driverBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	for i, v := range keys {
		log.Printf("driverBatcher.get: i:%v v:%v\n", i, v)
	}
	log.Printf("dataloader.driverBatcher.get, events: [%s]\n", strings.Join(keys.Keys(), ","))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// collect the keys to search for
	var teams []string
	for ix, key := range keys {
		id := key.String()
		teams = append(teams, id)
		keyOrder[id] = ix
	}
	// search for those users

	dbRecords := t.db.CollectDriversInTeams(ctx, teams)
	// if DB error, return
	// if err != nil {
	// 	return []*dataloader.Result{{Data: nil, Error: err}}
	// }
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for team, record := range dbRecords {
		ix, ok := keyOrder[team]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, team)
		}
	}
	// fill array positions with errors where not found in DB
	for team, ix := range keyOrder {
		err := fmt.Errorf("team not found %s", team)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}

// teams

// teamBatcher wraps storage and provides a "get" method for the driver dataloader
type teamBatcher struct {
	db storage.Storage
}

// get implements the dataloader for finding many tracks by Id and returns
// them in the order requested
func (t *teamBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	for i, v := range keys {
		log.Printf("teamBatcher: i:%v v:%v\n", i, v)
	}
	log.Printf("dataloader.teamBatcher.get, events: [%s]\n", strings.Join(keys.Keys(), ","))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// collect the keys to search for
	var teams []string
	for ix, key := range keys {
		id := key.String()
		teams = append(teams, id)
		keyOrder[id] = ix
	}
	// search for those users

	dbRecords := t.db.CollectTeamsForDrivers(ctx, teams)
	// if DB error, return
	// if err != nil {
	// 	return []*dataloader.Result{{Data: nil, Error: err}}
	// }
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for team, record := range dbRecords {
		ix, ok := keyOrder[team]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, team)
		}
	}
	// fill array positions with errors where not found in DB
	for team, ix := range keyOrder {
		err := fmt.Errorf("driver not found %s", team)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}
