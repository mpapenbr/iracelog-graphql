package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
	"github.com/mpapenbr/iracelog-graphql/internal/tracks"
)

type DbStorage struct {
	Storage
	pool *pgxpool.Pool
}

func NewDbStorage() *DbStorage {
	return &DbStorage{pool: database.InitDB()}
}

// tracks
func (db *DbStorage) GetAllTracks(ctx context.Context) ([]*model.Track, error) {

	var result []*model.Track

	tracks, err := tracks.GetALl(db.pool)
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, track := range tracks {
			result = append(result, &model.Track{ID: track.ID, Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length})
		}
	}
	return result, err
}
func (db *DbStorage) GetTracks(ctx context.Context, ids []int) ([]*model.Track, error) {

	var result []*model.Track

	tracks, err := tracks.GetByIds(db.pool, ids)
	// log.Printf("Tracks: %v\n", tracks)
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, track := range tracks {
			result = append(result, &model.Track{ID: track.ID, Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length})
		}
	}
	return result, err
}

// events
func (db *DbStorage) GetAllEvents(ctx context.Context) ([]*model.Event, error) {

	var result []*model.Event

	events, err := events.GetALl(db.pool)
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, event := range events {
			// this would cause assigning the last loop content to all result entries
			dbData := event
			result = append(result, &model.Event{ID: event.ID, Name: event.Name, Key: event.Key, TrackId: event.Info.TrackId, DbEvent: &dbData})
		}
	}
	return result, err
}

func (db *DbStorage) GetEvents(ctx context.Context, ids []int) ([]*model.Event, error) {

	var result []*model.Event

	events, err := events.GetByIds(db.pool, ids)
	// log.Printf("Events: %v\n", events)
	if err == nil {
		// convert the internal database Event to the GraphQL-Event
		for _, event := range events {
			// note: this is required. don't use the loop directly for DbEvent:&event.
			// this would cause assigning the last loop content to all result entries
			dbData := event
			result = append(result, &model.Event{ID: event.ID, Name: event.Name, Key: event.Key, TrackId: event.Info.TrackId, DbEvent: &dbData})
		}
	}
	return result, err
}

func (db *DbStorage) GetEventIdsForTrackId(ctx context.Context, trackId int) ([]int, error) {
	return events.GetIdsByTrackId(db.pool, trackId)
}

func (db *DbStorage) GetTeamsForEvent(ctx context.Context, event *model.Event) []*model.EventTeam {
	if event.DbEvent.Info.TeamRacing == 0 {
		return []*model.EventTeam{}
	}
	if event.DbAnalysisData == nil {
		log.Printf("loading analysis data for event %d\n", event.ID)
		x, _ := analysis.GetAnalysisForEvent(db.pool, event.ID)
		event.DbAnalysisData = x

	}
	ret := make([]*model.EventTeam, len(event.DbAnalysisData.CarInfo))
	for i, ci := range event.DbAnalysisData.CarInfo {
		drivers := make([]*model.Driver, len(ci.Drivers))
		for j, driver := range ci.Drivers {
			drivers[j] = &model.Driver{Name: driver.DriverName}
		}
		ret[i] = &model.EventTeam{Name: ci.Name, CarNum: ci.CarNum}

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

func (db *DbStorage) CollectDriversInTeams(ctx context.Context, teams []string) map[string][]*model.Driver {
	res, _ := analysis.SearchDriversInTeams(db.pool, teams)
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

func (db *DbStorage) CollectTeamsForDrivers(ctx context.Context, drivers []string) map[string][]*model.Team {
	res, _ := analysis.SearchTeamsForDrivers(db.pool, drivers)
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
