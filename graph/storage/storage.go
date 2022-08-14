package storage

import (
	"context"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
)

// This interface defines all available calls to storage.
type Storage interface {

	// GetTracks accepts many user IDs and returns an array of matching Tracks
	GetTracks(ctx context.Context, ids []int) ([]*model.Track, error)
	// GetEvents accepts many event IDs and returns an array of matching Tracks
	GetEvents(ctx context.Context, ids []int) ([]*model.Event, error)
	// GetEventIdsForTrack
	GetEventIdsForTrackId(ctx context.Context, trackId int) ([]int, error)

	// GetAllTracks lists all Tracks in the database
	GetAllTracks(ctx context.Context) ([]*model.Track, error)
	// GetAllEvents lists all Events in the database
	GetAllEvents(ctx context.Context) ([]*model.Event, error)

	// search drivers by name
	SearchDrivers(ctx context.Context, arg string) []*model.Driver
	// collect drivers for a given team name accross all events. returned map key is the team name
	CollectDriversInTeams(ctx context.Context, teams []string) map[string][]*model.Driver
	// collect drivers for a given team name accross all events. returned map key is the team name
	CollectTeamsForDrivers(ctx context.Context, drivers []string) map[string][]*model.Team
	// collect the eventIds for a specific driver (name)
	CollectEventIdsForDriver(ctx context.Context, driver string) []int
	// collect the eventIds for a specific team (name)
	CollectEventIdsForTeam(ctx context.Context, team string) []int
	// collect the analysis data for a number of eventIds
	CollectAnalysisData(ctx context.Context, eventIds []int) []analysis.DbAnalysis

	CollectEventIdsForTeams(ctx context.Context, teams []string) map[string][]int
	CollectEventIdsForDrivers(ctx context.Context, drivers []string) map[string][]int

	// search team by name
	SearchTeams(ctx context.Context, arg string) []*model.Team
}
