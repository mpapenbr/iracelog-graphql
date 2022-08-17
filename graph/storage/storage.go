package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
)

/*
This interface defines all available calls to storage.
A rule of thumb: for functions containing dataloader.Keys as argument the resulting map key is built by Key.String()
*/
type Storage interface {

	// GetTracks expects keys of type IntKey. IntKey.String() is used as map key
	GetTracksByKeys(ctx context.Context, keys dataloader.Keys) map[string]*model.Track
	// GetEvents expects keys of type IntKey. IntKey.String() is used as map key
	GetEventsByKeys(ctx context.Context, keys dataloader.Keys) map[string]*model.Event

	// trackIds contains IntKey instances.
	GetEventsForTrackIdsKeys(ctx context.Context, trackIds dataloader.Keys) map[string][]*model.Event

	// GetAllTracks lists all Tracks in the database
	GetAllTracks(ctx context.Context) ([]*model.Track, error)
	// GetAllEvents lists all Events in the database
	GetAllEvents(ctx context.Context) ([]*model.Event, error)

	// search drivers by name
	SearchDrivers(ctx context.Context, arg string) []*model.Driver
	// collect drivers for given team name (StringKey) accross all events. returned map key is the team name
	CollectDriversInTeams(ctx context.Context, teams dataloader.Keys) map[string][]*model.Driver
	// collect teams for a given driver name (StringKey) accross all events. returned map key is the driver name
	CollectTeamsForDrivers(ctx context.Context, drivers dataloader.Keys) map[string][]*model.Team

	// collect the analysis data for a number of eventIds
	CollectAnalysisData(ctx context.Context, eventIds dataloader.Keys) map[string]analysis.DbAnalysis

	CollectEventIdsForTeams(ctx context.Context, teams dataloader.Keys) map[string][]int
	CollectEventIdsForDrivers(ctx context.Context, drivers dataloader.Keys) map[string][]int

	// search team by name
	SearchTeams(ctx context.Context, arg string) []*model.Team
}
