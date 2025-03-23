package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
)

// This interface defines all available calls to storage.
// A rule of thumb: for functions containing dataloader.Keys as argument
// the resulting map key is built by Key.String()

type Storage interface {
	ResolveTenant(ctx context.Context, externalId string) (int, error)
	// GetTracks expects keys of type IntKey. IntKey.String() is used as map key
	GetTracksByKeys(ctx context.Context, keys dataloader.Keys) map[string]*model.Track
	// GetEvents expects keys of type IntKey. IntKey.String() is used as map key
	GetEventsByKeys(ctx context.Context, keys dataloader.Keys) map[string]*model.Event

	// trackIds contains IntKey instances.
	GetEventsForTrackIdsKeys(
		ctx context.Context,
		trackIds dataloader.Keys) map[string][]*model.Event

	// GetAllTracks lists all Tracks in the database
	GetAllTracks(
		ctx context.Context,
		limit *int,
		offset *int,
		sort []*model.TrackSortArg) ([]*model.Track, error)
	// GetAllEvents lists all Events in the database
	GetAllEvents(
		ctx context.Context,
		limit *int,
		offset *int,
		sort []*model.EventSortArg) ([]*model.Event, error)

	// simple search events by name,description,driver.name,team.name,car.name,track.name
	SimpleSearchEvents(
		ctx context.Context,
		arg string,
		limit *int,
		offset *int,
		sort []*model.EventSortArg) ([]*model.Event, error)
	// advanced search events.
	// arg is examined for search keys (like name,track,driver,team,car)
	AdvancedSearchEvents(
		ctx context.Context,
		arg *events.EventSearchKeys,
		limit *int,
		offset *int,
		sort []*model.EventSortArg) ([]*model.Event, error)

	// search drivers by name
	SearchDrivers(ctx context.Context, arg string) []*model.Driver
	// collect drivers for given team name (StringKey) across all events.
	// returned map key is the team name
	CollectDriversInTeams(
		ctx context.Context,
		teams dataloader.Keys) map[string][]*model.Driver
	// collect teams for a given driver name (StringKey) across all events.
	// returned map key is the driver name
	CollectTeamsForDrivers(
		ctx context.Context,
		drivers dataloader.Keys) map[string][]*model.Team

	// collect the analysis data for a number of eventIds
	CollectAnalysisData(
		ctx context.Context,
		eventIds dataloader.Keys) map[string]analysis.DbAnalysis

	CollectEventIdsForTeams(
		ctx context.Context,
		teams dataloader.Keys) map[string][]int
	CollectEventIdsForDrivers(
		ctx context.Context,
		drivers dataloader.Keys) map[string][]int

	// collect the drivers for a number of eventIds
	CollectEventDrivers(
		ctx context.Context,
		eventIds dataloader.Keys) map[string][]*model.EventDriver

	// search team by name
	SearchTeams(ctx context.Context, arg string) []*model.Team

	// new collectors start here
	// collect the event entries for a number of eventIds
	CollectEventEntries(
		ctx context.Context,
		eventIds dataloader.Keys) map[string][]*model.EventEntry
	// collect the event entries for selected ids
	CollectEventEntriesById(
		ctx context.Context,
		ids dataloader.Keys) map[string]*model.EventEntry
	// collect the cars for a number of eventIds
	CollectEventCars(
		ctx context.Context,
		eventIds dataloader.Keys) map[string][]*model.Car

	// collect the cars for a number of eventEntryIds
	CollectCarsByEventEntry(
		ctx context.Context,
		eventEntryIds dataloader.Keys) map[string]*model.Car
	// collect the teams for a number of eventEntryIds
	CollectTeamByEventEntry(
		ctx context.Context,
		eventEntryIds dataloader.Keys) map[string]*model.EventTeam
	// collect the event drivers for a number of eventEntryIds
	CollectDriversByEventEntry(
		ctx context.Context,
		eventEntryIds dataloader.Keys) map[string][]*model.EventDriver
	// collect the event drivers for a number of eventEntryIds
	CollectDriversByTeam(
		ctx context.Context,
		teamIds dataloader.Keys) map[string][]*model.EventDriver
}
