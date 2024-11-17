package dataloader

import (
	"context"
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
	trackLoader               *dataloader.Loader
	eventLoader               *dataloader.Loader
	driverLoader              *dataloader.Loader
	teamLoader                *dataloader.Loader
	analysisLoader            *dataloader.Loader // deprecated
	teamEventLinkLoader       *dataloader.Loader // deprecated
	driverEventLinkLoader     *dataloader.Loader // deprecated
	eventsByTrackLoader       *dataloader.Loader // deprecated
	driversByEventLoader      *dataloader.Loader // deprecated
	eventEntriesByIdsLoader   *dataloader.Loader
	eventEntriesByEventLoader *dataloader.Loader
	eventCarsLoader           *dataloader.Loader
	eventEntryCarLoader       *dataloader.Loader
	driversByEventEntryLoader *dataloader.Loader
	driversByTeamLoader       *dataloader.Loader
	teamByEventEntryLoader    *dataloader.Loader
}

// deprecated
func (i *DataLoader) GetTeamDrivers(ctx context.Context, team string) ([]*model.Driver, []error) {
	thunk := i.driverLoader.Load(ctx, gopher_dataloader.StringKey(team))
	result, err := thunk()
	if err != nil {
		return nil, nil
	}
	return result.([]*model.Driver), nil
}

// deprecated
func (i *DataLoader) GetDriversTeams(ctx context.Context, driver string) ([]*model.Team, []error) {
	thunk := i.teamLoader.Load(ctx, gopher_dataloader.StringKey(driver))
	result, err := thunk()
	if err != nil {
		return nil, nil
	}
	return result.([]*model.Team), nil
}

// deprecated
func (i *DataLoader) GetEventIdsForTeam(ctx context.Context, team string) []int {
	thunk := i.teamEventLinkLoader.Load(ctx, gopher_dataloader.StringKey(team))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil
	}
	return result.([]int)
}

// deprecated
func (i *DataLoader) GetEventIdsForDriver(ctx context.Context, driver string) []int {
	thunk := i.driverEventLinkLoader.Load(ctx, gopher_dataloader.StringKey(driver))
	result, err := thunk()
	if err != nil {
		log.Printf("error loading analysis data: %v", err)
		return nil
	}
	return result.([]int)
}

// deprecated

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
	driversByEvent := &genericMapBatcher[[]*model.EventDriver]{collector: db.CollectEventDrivers}

	// new collectors start here
	eventEntriesByIds := &genericMapBatcher[*model.EventEntry]{collector: db.CollectEventEntriesById}
	eventEntriesByEvent := &genericMapBatcher[[]*model.EventEntry]{collector: db.CollectEventEntries}
	eventEntryCars := &genericMapBatcher[*model.Car]{collector: db.CollectCarsByEventEntry}
	eventCars := &genericMapBatcher[[]*model.Car]{collector: db.CollectEventCars}
	driversByEventEntry := &genericMapBatcher[[]*model.EventDriver]{collector: db.CollectDriversByEventEntry}
	driversByTeam := &genericMapBatcher[[]*model.EventDriver]{collector: db.CollectDriversByTeam}
	teamByEventEntry := &genericMapBatcher[*model.EventTeam]{collector: db.CollectTeamByEventEntry}

	loaders := &DataLoader{
		trackLoader:               dataloader.NewBatchedLoader(tracks.get),
		eventLoader:               dataloader.NewBatchedLoader(events.get),
		driverLoader:              dataloader.NewBatchedLoader(drivers.get),
		teamLoader:                dataloader.NewBatchedLoader(teams.get),
		analysisLoader:            dataloader.NewBatchedLoader(analysis.get),
		teamEventLinkLoader:       dataloader.NewBatchedLoader(teamEventLink.get),
		driverEventLinkLoader:     dataloader.NewBatchedLoader(driverEventLink.get),
		eventsByTrackLoader:       dataloader.NewBatchedLoader(eventsByTrack.get),
		driversByEventLoader:      dataloader.NewBatchedLoader(driversByEvent.get),
		eventEntriesByIdsLoader:   dataloader.NewBatchedLoader(eventEntriesByIds.get),
		eventEntriesByEventLoader: dataloader.NewBatchedLoader(eventEntriesByEvent.get),
		eventCarsLoader:           dataloader.NewBatchedLoader(eventCars.get),
		eventEntryCarLoader:       dataloader.NewBatchedLoader(eventEntryCars.get),
		driversByEventEntryLoader: dataloader.NewBatchedLoader(driversByEventEntry.get),
		driversByTeamLoader:       dataloader.NewBatchedLoader(driversByTeam.get),
		teamByEventEntryLoader:    dataloader.NewBatchedLoader(teamByEventEntry.get),
	}
	return loaders
}
