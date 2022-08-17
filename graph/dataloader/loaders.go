package dataloader

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
	genericLoader         *dataloader.Loader
}

// GetTrack wraps the Track dataloader for efficient retrieval by track ID
func (i *DataLoader) GetTrack(ctx context.Context, trackID int) (*model.Track, error) {
	thunk := i.trackLoader.Load(ctx, storage.IntKey(trackID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.Track), nil
}

func (i *DataLoader) GetTracks(ctx context.Context, trackIds []int) ([]*model.Track, []error) {

	trackKeys := storage.NewKeysFromInts(trackIds)
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
	eventKeys := storage.NewKeysFromInts(eventIDs)
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
	thunk := i.eventLoader.Load(ctx, storage.IntKey(eventID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.Event), nil
}

func (i *DataLoader) GetEventsForTrack(ctx context.Context, trackId int) []*model.Event {
	// thunk := i.eventsByTrackLoader.Load(ctx, gopher_dataloader.StringKey(fmt.Sprintf("%d", trackId)))
	thunk := i.eventsByTrackLoader.Load(ctx, storage.IntKey(trackId))
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

	tracks := &genericMapBatcher[*model.Track]{collector: db.GetTracksByKeys}
	events := &genericMapBatcher[*model.Event]{collector: db.GetEventsByKeys}

	drivers := &genericMapBatcher[[]*model.Driver]{collector: db.CollectDriversInTeams}
	teams := &genericMapBatcher[[]*model.Team]{collector: db.CollectTeamsForDrivers}

	analysis := &genericMapBatcher[analysis.DbAnalysis]{collector: db.CollectAnalysisData}

	driverEventLink := &genericMapBatcher[[]int]{collector: db.CollectEventIdsForDrivers}
	teamEventLink := &genericMapBatcher[[]int]{collector: db.CollectEventIdsForTeams}

	eventsByTrack := &genericMapBatcher[[]*model.Event]{collector: db.GetEventsForTrackIdsKeys}

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
