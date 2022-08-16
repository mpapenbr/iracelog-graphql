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
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
)

type ctxKey string

const loadersKey = ctxKey("dataloader")

type DataLoader struct {
	trackLoader           *dataloader.Loader
	eventLoader           *dataloader.Loader
	driverLoader          *dataloader.Loader
	teamLoader            *dataloader.Loader
	analysisLoader        *dataloader.Loader
	teamEventLinkLoader   *dataloader.Loader
	driverEventLinkLoader *dataloader.Loader
	eventsByTrackLoader   *dataloader.Loader
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

func (i *DataLoader) GetEventsForTrack(ctx context.Context, trackId int) []*model.Event {
	thunk := i.eventsByTrackLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", trackId)))
	result, err := thunk()
	if err != nil {
		return nil
	}
	return result.([]*model.Event)
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

func (i *DataLoader) GetEventTeams(ctx context.Context, eventId int) ([]*model.EventTeam, []error) {
	thunk := i.analysisLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", eventId)))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil, nil
	}
	ret := []*model.EventTeam{}

	dbData := result.(analysis.DbAnalysis)
	for _, ci := range dbData.CarInfo {
		drivers := make([]*model.EventDriver, len(ci.Drivers))
		for j, driver := range ci.Drivers {
			drivers[j] = &model.EventDriver{Name: driver.DriverName}
		}
		ret = append(ret, &model.EventTeam{Name: ci.Name, CarNum: ci.CarNum, Drivers: drivers})
	}
	return ret, nil
}

func (i *DataLoader) GetEventIdsForTeam(ctx context.Context, team string) []int {
	thunk := i.teamEventLinkLoader.Load(ctx, gopher_dataloader.StringKey(team))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil
	}
	return result.([]int)
}

func (i *DataLoader) GetEventIdsForDriver(ctx context.Context, driver string) []int {
	thunk := i.driverEventLinkLoader.Load(ctx, gopher_dataloader.StringKey(driver))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil
	}
	return result.([]int)
}

func (i *DataLoader) GetEventDrivers(ctx context.Context, eventId int) ([]*model.EventDriver, []error) {
	thunk := i.analysisLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", eventId)))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil, nil
	}
	ret := []*model.EventDriver{}

	dbData := result.(analysis.DbAnalysis)
	for _, ci := range dbData.CarInfo {
		for _, driver := range ci.Drivers {
			ret = append(ret, &model.EventDriver{Name: driver.DriverName})
		}
	}
	return ret, nil
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
	tracks := &genericByIdBatcher[model.Track]{db: db, collector: db.GetTracks, idExtractor: func(entity *model.Track) int { return entity.ID }}
	events := &genericByIdBatcher[model.Event]{db: db, collector: db.GetEvents, idExtractor: func(entity *model.Event) int { return entity.ID }}

	// once we have a "real" generic matcher this should be converted to use map[int]...
	eventsByTrack := &genericByNameBatcher[model.Event]{db: db, collector: db.GetEventsForTrackIds}
	drivers := &genericByNameBatcher[model.Driver]{db: db, collector: db.CollectDriversInTeams}
	teams := &genericByNameBatcher[model.Team]{db: db, collector: db.CollectTeamsForDrivers}

	analysis := &analysisBatcher{db: db}
	driverEventLink := &genericEventLinkBatcher{db: db, collector: db.CollectEventIdsForDrivers}
	teamEventLink := &genericEventLinkBatcher{db: db, collector: db.CollectEventIdsForTeams}
	loaders := &DataLoader{
		trackLoader:           dataloader.NewBatchedLoader(tracks.get),
		eventLoader:           dataloader.NewBatchedLoader(events.get),
		driverLoader:          dataloader.NewBatchedLoader(drivers.get),
		teamLoader:            dataloader.NewBatchedLoader(teams.get),
		analysisLoader:        dataloader.NewBatchedLoader(analysis.get),
		teamEventLinkLoader:   dataloader.NewBatchedLoader(teamEventLink.get),
		driverEventLinkLoader: dataloader.NewBatchedLoader(driverEventLink.get),
		eventsByTrackLoader:   dataloader.NewBatchedLoader(eventsByTrack.get),
	}
	return loaders
}

type ByNameCollector[E any] func(ctx context.Context, names []string) map[string][]*E

type genericByNameBatcher[E any] struct {
	db        storage.Storage
	collector ByNameCollector[E]
}

// get implements the dataloader for finding many tracks by Id and returns
// them in the order requested
func (t *genericByNameBatcher[E]) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// for i, v := range keys {
	// 	log.Printf("driverBatcher.get: i:%v v:%v\n", i, v)
	// }
	log.Printf("dataloader.genericByNameBatcher.get, names: [%s]\n", strings.Join(keys.Keys(), ","))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// collect the keys to search for
	var refNames []string
	for ix, key := range keys {
		id := key.String()
		refNames = append(refNames, id)
		keyOrder[id] = ix
	}
	// search for those users

	dbRecords := t.collector(ctx, refNames)
	// if DB error, return
	// if err != nil {
	// 	return []*dataloader.Result{{Data: nil, Error: err}}
	// }
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for refName, record := range dbRecords {
		ix, ok := keyOrder[refName]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, refName)
		}
	}
	// fill array positions with errors where not found in DB
	for refName, ix := range keyOrder {
		err := fmt.Errorf("refName not found %s", refName)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}

// analysisBatcher wraps storage and provides a "get" method for the analysis dataloader
type analysisBatcher struct {
	db storage.Storage
}

// get implements the dataloader for finding many tracks by Id and returns
// them in the order requested
func (t *analysisBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	for i, v := range keys {
		log.Printf("analysisBatcher: i:%v v:%v\n", i, v)
	}
	log.Printf("dataloader.analysisBatcher.get, events: [%s]\n", strings.Join(keys.Keys(), ","))
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

	dbRecords := t.db.CollectAnalysisData(ctx, eventIDs)
	// if DB error, return
	// if err != nil {
	// 	return []*dataloader.Result{{Data: nil, Error: err}}
	// }
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for _, record := range dbRecords {
		ix, ok := keyOrder[record.EventId]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, record.EventId)
		}
	}
	// fill array positions with errors where not found in DB
	for eventId, ix := range keyOrder {
		err := fmt.Errorf("event not found %d", eventId)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}

type ByIdCollector[E any] func(ctx context.Context, ids []int) ([]*E, error)
type IntIdExtractor[E any] func(entity E) int

type genericByIdBatcher[E any] struct {
	db          storage.Storage
	collector   ByIdCollector[E]
	idExtractor IntIdExtractor[*E]
}

func (t *genericByIdBatcher[E]) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// for i, v := range keys {
	// 	log.Printf("i:%v v:%v\n", i, v)
	// }
	log.Printf("dataloader.genericByIdBatcher.get, ids: [%s]\n", strings.Join(keys.Keys(), ","))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[int]int, len(keys))
	// collect the keys to search for
	var ids []int
	for ix, key := range keys {
		id, _ := strconv.Atoi(key.String())
		ids = append(ids, id)
		keyOrder[id] = ix
	}
	// search for those users

	dbRecords, err := t.collector(ctx, ids)
	// if DB error, return
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for _, record := range dbRecords {
		ix, ok := keyOrder[t.idExtractor(record)]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, t.idExtractor(record))
		}
	}
	// fill array positions with errors where not found in DB
	for id, ix := range keyOrder {
		err := fmt.Errorf("item not found %d", id)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}

type IdCollector func(ctx context.Context, names []string) map[string][]int
type genericEventLinkBatcher struct {
	db        storage.Storage
	collector IdCollector
}

// get implements the dataloader for finding many tracks by Id and returns
// them in the order requested
func (t *genericEventLinkBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// for i, v := range keys {
	// 	log.Printf("genericEventLinkBatcher: i:%v v:%v\n", i, v)
	// }
	log.Printf("dataloader.genericEventLinkBatcher.get, for names: [%s]\n", strings.Join(keys.Keys(), ","))
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// collect the keys to search for
	var drivers []string
	for ix, key := range keys {
		id := key.String()
		drivers = append(drivers, id)
		keyOrder[id] = ix
	}
	// search for those users

	dbRecords := t.collector(ctx, drivers)
	// if DB error, return
	// if err != nil {
	// 	return []*dataloader.Result{{Data: nil, Error: err}}
	// }
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// enumerate records, put into output
	for driver, record := range dbRecords {
		ix, ok := keyOrder[driver]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, driver)
		}
	}
	// fill array positions with errors where not found in DB
	for driver, ix := range keyOrder {
		err := fmt.Errorf("name not found %s", driver)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	// return results
	return results
}
