package storage

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
	"github.com/mpapenbr/iracelog-graphql/internal/tracks"
)

type DbStorage struct {
	// Storage
	pool *pgxpool.Pool
}

func NewDbStorage() *DbStorage {
	return &DbStorage{pool: database.InitDB()}
}

func NewDbStorageWithPool(pool *pgxpool.Pool) *DbStorage {
	return &DbStorage{pool: pool}
}

// tracks
func (db *DbStorage) GetAllTracks(ctx context.Context, limit *int, offset *int, sort []*model.TrackSortArg) ([]*model.Track, error) {
	var result []*model.Track

	dbTrackSortArg := convertTrackSortArgs(sort)
	tracks, err := tracks.GetALl(db.pool, internal.DbPageable{Limit: limit, Offset: offset, Sort: dbTrackSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, track := range tracks {
			result = append(result, convertDbTrackToModel(track))
		}
	}
	return result, err
}

func (db *DbStorage) GetTracksByKeys(ctx context.Context, ids dataloader.Keys) map[string]*model.Track {
	intIds := IntKeysToSlice(ids)
	result := map[string]*model.Track{}

	tracks, _ := tracks.GetByIds(db.pool, intIds)
	// log.Printf("Tracks: %v\n", tracks)

	// convert the internal database Track to the GraphQL-Track
	for _, track := range tracks {
		result[IntKey(track.ID).String()] = convertDbTrackToModel(track)
	}

	return result
}

// events
func (db *DbStorage) GetAllEvents(ctx context.Context, limit *int, offset *int, sort []*model.EventSortArg) ([]*model.Event, error) {
	var result []*model.Event
	dbEventSortArg := convertEventSortArgs(sort)
	events, err := events.GetALl(db.pool, internal.DbPageable{Limit: limit, Offset: offset, Sort: dbEventSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, dbEvents := range events {
			// this would cause assigning the last loop content to all result entries

			result = append(result, convertDbEventToModel(dbEvents))
		}
	}
	return result, err
}

func (db *DbStorage) SimpleSearchEvents(ctx context.Context, arg string, limit *int, offset *int, sort []*model.EventSortArg) ([]*model.Event, error) {
	var result []*model.Event
	dbEventSortArg := convertEventSortArgs(sort)
	events, err := events.SimpleEventSearch(db.pool, arg, internal.DbPageable{Limit: limit, Offset: offset, Sort: dbEventSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, dbEvents := range events {
			// this would cause assigning the last loop content to all result entries

			result = append(result, convertDbEventToModel(dbEvents))
		}
	}
	return result, err
}

func (db *DbStorage) AdvancedSearchEvents(ctx context.Context, arg *events.EventSearchKeys, limit *int, offset *int, sort []*model.EventSortArg) ([]*model.Event, error) {
	var result []*model.Event
	dbEventSortArg := convertEventSortArgs(sort)
	events, err := events.AdvancedEventSearch(db.pool, arg, internal.DbPageable{Limit: limit, Offset: offset, Sort: dbEventSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, dbEvents := range events {
			// this would cause assigning the last loop content to all result entries

			result = append(result, convertDbEventToModel(dbEvents))
		}
	}
	return result, err
}

func (db *DbStorage) GetEventsByKeys(ctx context.Context, ids dataloader.Keys) map[string]*model.Event {
	intIds := IntKeysToSlice(ids)
	result := map[string]*model.Event{}

	events, _ := events.GetByIds(db.pool, intIds)
	// log.Printf("Tracks: %v\n", tracks)

	// convert the internal database Track to the GraphQL-Track
	for _, dbEvents := range events {
		// this would cause assigning the last loop content to all result entries
		result[IntKey(dbEvents.ID).String()] = convertDbEventToModel(dbEvents)
	}

	return result
}

// Note: we use (temporary) a string as key (to reuse existing batcher mechanics)
func (db *DbStorage) GetEventsForTrackIdsKeys(ctx context.Context, trackIds dataloader.Keys) map[string][]*model.Event {
	result := map[string][]*model.Event{}

	intTrackIds := make([]int, len(trackIds))
	for i, id := range trackIds {
		intTrackIds[i] = id.Raw().(int)
	}
	byTrackId, err := events.GetEventsByTrackIds(db.pool, intTrackIds, internal.DbPageable{Sort: convertEventSortArgs([]*model.EventSortArg{})})
	// log.Printf("Events: %v\n", events)
	if err == nil {
		// convert the internal database Event to the GraphQL-Event
		for k, event := range byTrackId {
			convertedEvents := make([]*model.Event, len(event))
			for i, dbEvent := range event {
				convertedEvents[i] = convertDbEventToModel(*dbEvent)
			}
			result[fmt.Sprintf("%d", k)] = convertedEvents
		}
	}
	return result
}

func (db *DbStorage) CollectAnalysisData(ctx context.Context, eventIds dataloader.Keys) map[string]analysis.DbAnalysis {
	res, _ := analysis.GetAnalysisForEvents(db.pool, IntKeysToSlice(eventIds))
	ret := map[string]analysis.DbAnalysis{}
	for _, a := range res {
		ret[IntKey(a.EventId).String()] = a
	}
	return ret
}

func (db *DbStorage) SearchDrivers(ctx context.Context, arg string) []*model.Driver {
	res, _ := analysis.SearchDrivers(db.pool, arg)
	ret := make([]*model.Driver, len(res))
	for i, d := range res {
		teams := make([]*model.Team, len(d.Teams))
		for j, d := range d.Teams {
			teams[j] = &model.Team{Name: d}
		}
		ret[i] = &model.Driver{Name: d.Name, Teams: teams, CarNum: d.CarNum, CarClass: d.CarClass}
	}
	return ret
}

func (db *DbStorage) CollectDriversInTeams(ctx context.Context, teams dataloader.Keys) map[string][]*model.Driver {
	res, _ := analysis.SearchDriversInTeams(db.pool, teams.Keys())
	ret := map[string][]*model.Driver{}
	for k, v := range res {
		drivers := make([]*model.Driver, len(v))
		for i, d := range v {
			drivers[i] = &model.Driver{Name: d.Name, CarNum: d.CarNum, CarClass: d.CarClass}
		}
		ret[k] = drivers
	}
	return ret
}

func (db *DbStorage) CollectTeamsForDrivers(ctx context.Context, drivers dataloader.Keys) map[string][]*model.Team {
	res, _ := analysis.SearchTeamsForDrivers(db.pool, drivers.Keys())
	ret := map[string][]*model.Team{}
	for k, v := range res {
		teams := make([]*model.Team, len(v))
		for i, d := range v {
			teams[i] = &model.Team{Name: d.Name, CarNum: d.CarNum, CarClass: d.CarClass}
		}
		ret[k] = teams
	}
	return ret
}

func (db *DbStorage) CollectEventIdsForTeams(ctx context.Context, teams dataloader.Keys) map[string][]int {
	ret, _ := analysis.CollectEventIdsForTeams(db.pool, teams.Keys())
	return ret
}

func (db *DbStorage) CollectEventIdsForDrivers(ctx context.Context, drivers dataloader.Keys) map[string][]int {
	ret, _ := analysis.CollectEventIdsForDrivers(db.pool, drivers.Keys())
	return ret
}

func (db *DbStorage) SearchTeams(ctx context.Context, arg string) []*model.Team {
	res, _ := analysis.SearchTeams(db.pool, arg)
	ret := make([]*model.Team, len(res))
	for i, d := range res {
		drivers := make([]*model.Driver, len(d.Drivers))
		for j, d := range d.Drivers {
			drivers[j] = &model.Driver{Name: d}
		}
		ret[i] = &model.Team{Name: d.Name, Drivers: drivers, CarNum: d.CarNum, CarClass: d.CarClass}
	}
	return ret
}
