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

func (db *DbStorage) GetTeamsForEvent(ctx context.Context, event *model.Event) []*model.Team {
	if event.DbEvent.Info.TeamRacing == 0 {
		return []*model.Team{}
	}
	if event.DbAnalysisData == nil {
		log.Printf("loading analysis data for event %d\n", event.ID)
		x, _ := analysis.GetAnalysisForEvent(db.pool, event.ID)
		event.DbAnalysisData = x

	}
	ret := make([]*model.Team, len(event.DbAnalysisData.CarInfo))
	for i, ci := range event.DbAnalysisData.CarInfo {
		drivers := make([]*model.Driver, len(ci.Drivers))
		for j, driver := range ci.Drivers {
			drivers[j] = &model.Driver{Name: driver.DriverName, CarNum: ci.CarNum}
		}
		ret[i] = &model.Team{Name: ci.Name, CarNum: ci.CarNum, Drivers: drivers}

	}
	return ret
}
