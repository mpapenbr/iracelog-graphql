package dataloader

import (
	"context"
	"net/http"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/graph/storage"
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
	"github.com/mpapenbr/iracelog-graphql/log"
)

type ctxKey string

const loadersKey = ctxKey("dataloader")

type DataLoader struct {
	l                         *log.Logger
	trackLoader               *dataloader.Loader
	eventLoader               *dataloader.Loader
	driverLoader              *dataloader.Loader
	teamLoader                *dataloader.Loader
	analysisLoader            *dataloader.Loader // deprecated
	teamEventLinkLoader       *dataloader.Loader // deprecated
	driverEventLinkLoader     *dataloader.Loader // deprecated
	eventsByTrackLoader       *dataloader.Loader // deprecated
	driversByEventLoader      *dataloader.Loader // deprecated
	eventEntriesByIDsLoader   *dataloader.Loader
	eventEntriesByEventLoader *dataloader.Loader
	eventCarsLoader           *dataloader.Loader
	eventEntryCarLoader       *dataloader.Loader
	driversByEventEntryLoader *dataloader.Loader
	driversByTeamLoader       *dataloader.Loader
	teamByEventEntryLoader    *dataloader.Loader
}

// deprecated
//
//nolint:all // this is a deprecated function
func (i *DataLoader) GetTeamDrivers(
	ctx context.Context,
	team string,
) ([]*model.Driver, []error) {
	thunk := i.driverLoader.Load(ctx, dataloader.StringKey(team))
	result, err := thunk()
	if err != nil {
		return nil, nil
	}
	return result.([]*model.Driver), nil
}

// deprecated
//
//nolint:all // this is a deprecated function
func (i *DataLoader) GetDriversTeams(
	ctx context.Context,
	driver string,
) ([]*model.Team, []error) {
	thunk := i.teamLoader.Load(ctx, dataloader.StringKey(driver))
	result, err := thunk()
	if err != nil {
		return nil, nil
	}
	return result.([]*model.Team), nil
}

// deprecated
//
//nolint:all // this is a deprecated function
func (i *DataLoader) GetEventIDsForTeam(ctx context.Context, team string) []int {
	thunk := i.teamEventLinkLoader.Load(ctx, dataloader.StringKey(team))
	result, err := thunk()
	if err != nil {
		i.l.Error("error loading analysis data", log.ErrorField(err))
		return nil
	}
	return result.([]int)
}

// deprecated
//
//nolint:all // this is a deprecated function
func (i *DataLoader) GetEventIDsForDriver(ctx context.Context, driver string) []int {
	thunk := i.driverEventLinkLoader.Load(ctx, dataloader.StringKey(driver))
	result, err := thunk()
	if err != nil {
		i.l.Error("error loading analysis data", log.ErrorField(err))
		return nil
	}
	return result.([]int)
}

func Middleware(db storage.Storage, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := NewDataLoader(db) // we want fresh loaders on each request
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *DataLoader {
	//nolint:errcheck // we are sure that the type is correct
	return ctx.Value(loadersKey).(*DataLoader)
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
//
//nolint:lll,funlen // by design
func NewDataLoader(db storage.Storage) *DataLoader {
	// define the data loader

	tracks := &genericMapBatcher[*model.Track]{
		name: "trackByKeys", collector: db.GetTracksByKeys,
	}
	events := &genericMapBatcher[*model.Event]{
		name: "eventsByKeys", collector: db.GetEventsByKeys,
	}

	drivers := &genericMapBatcher[[]*model.Driver]{
		name: "driversInTeams", collector: db.CollectDriversInTeams,
	}
	teams := &genericMapBatcher[[]*model.Team]{
		name: "teamsForDrivers", collector: db.CollectTeamsForDrivers,
	}

	analysisData := &genericMapBatcher[analysis.DBAnalysis]{collector: db.CollectAnalysisData}

	driverEventLink := &genericMapBatcher[[]int]{
		name: "eventIDsForDriver", collector: db.CollectEventIDsForDrivers,
	}
	teamEventLink := &genericMapBatcher[[]int]{
		name: "eventIDsForTeams", collector: db.CollectEventIDsForTeams,
	}

	eventsByTrack := &genericMapBatcher[[]*model.Event]{
		name: "eventsByTrack", collector: db.GetEventsForTrackIDsKeys,
	}
	driversByEvent := &genericMapBatcher[[]*model.EventDriver]{
		name: "eventDrivers", collector: db.CollectEventDrivers,
	}

	// new collectors start here
	eventEntriesByIDs := &genericMapBatcher[*model.EventEntry]{
		name: "eventEntriesByID", collector: db.CollectEventEntriesByID,
	}
	eventEntriesByEvent := &genericMapBatcher[[]*model.EventEntry]{
		name: "eventEntries", collector: db.CollectEventEntries,
	}
	eventEntryCars := &genericMapBatcher[*model.Car]{
		name: "carsByEventEntry", collector: db.CollectCarsByEventEntry,
	}
	eventCars := &genericMapBatcher[[]*model.Car]{
		name: "eventCars", collector: db.CollectEventCars,
	}
	driversByEventEntry := &genericMapBatcher[[]*model.EventDriver]{
		name: "driversByEventEntry", collector: db.CollectDriversByEventEntry,
	}
	driversByTeam := &genericMapBatcher[[]*model.EventDriver]{
		name: "driversByTeam", collector: db.CollectDriversByTeam,
	}
	teamByEventEntry := &genericMapBatcher[*model.EventTeam]{
		name: "teamByEventEntry", collector: db.CollectTeamByEventEntry,
	}

	loaders := &DataLoader{
		trackLoader:               dataloader.NewBatchedLoader(tracks.get),
		eventLoader:               dataloader.NewBatchedLoader(events.get),
		driverLoader:              dataloader.NewBatchedLoader(drivers.get),
		teamLoader:                dataloader.NewBatchedLoader(teams.get),
		analysisLoader:            dataloader.NewBatchedLoader(analysisData.get),
		teamEventLinkLoader:       dataloader.NewBatchedLoader(teamEventLink.get),
		driverEventLinkLoader:     dataloader.NewBatchedLoader(driverEventLink.get),
		eventsByTrackLoader:       dataloader.NewBatchedLoader(eventsByTrack.get),
		driversByEventLoader:      dataloader.NewBatchedLoader(driversByEvent.get),
		eventEntriesByIDsLoader:   dataloader.NewBatchedLoader(eventEntriesByIDs.get),
		eventEntriesByEventLoader: dataloader.NewBatchedLoader(eventEntriesByEvent.get),
		eventCarsLoader:           dataloader.NewBatchedLoader(eventCars.get),
		eventEntryCarLoader:       dataloader.NewBatchedLoader(eventEntryCars.get),
		driversByEventEntryLoader: dataloader.NewBatchedLoader(driversByEventEntry.get),
		driversByTeamLoader:       dataloader.NewBatchedLoader(driversByTeam.get),
		teamByEventEntryLoader:    dataloader.NewBatchedLoader(teamByEventEntry.get),
	}
	return loaders
}
