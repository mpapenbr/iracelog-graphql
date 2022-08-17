package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/mpapenbr/iracelog-graphql/graph/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/generated"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
)

// Teams is the resolver for the teams field.
func (r *driverResolver) Teams(ctx context.Context, obj *model.Driver) ([]*model.Team, error) {
	dbResults, _ := dataloader.For(ctx).GetDriversTeams(ctx, obj.Name)
	// log.Printf("dbResult: %v err: %v\n", dbResults, err)
	return dbResults, nil
}

// Events is the resolver for the events field.
func (r *driverResolver) Events(ctx context.Context, obj *model.Driver) ([]*model.Event, error) {
	eventIds := dataloader.For(ctx).GetEventIdsForDriver(ctx, obj.Name)
	tmp, _ := dataloader.For(ctx).GetEvents(ctx, eventIds)
	return tmp, nil
}

// Track is the resolver for the track field.
func (r *eventResolver) Track(ctx context.Context, obj *model.Event) (*model.Track, error) {
	// track := tracks.GetById(obj.TrackId)
	// result := &model.Track{ID: track.ID, Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length}
	// fmt.Printf("eventResolver.Track, event=%v, track=%v\n", obj.ID, obj.TrackId)
	return dataloader.For(ctx).GetTrack(ctx, obj.TrackId)
}

// Teams is the resolver for the teams field.
func (r *eventResolver) Teams(ctx context.Context, obj *model.Event) ([]*model.EventTeam, error) {
	ret, _ := dataloader.For(ctx).GetEventTeams(ctx, obj.ID)
	return ret, nil
	// Old single version: return r.db.GetTeamsForEvent(ctx, obj), nil
}

// Drivers is the resolver for the drivers field.
func (r *eventResolver) Drivers(ctx context.Context, obj *model.Event) ([]*model.EventDriver, error) {
	ret, _ := dataloader.For(ctx).GetEventDrivers(ctx, obj.ID)
	return ret, nil
}

// GetEvents is the resolver for the getEvents field.
func (r *queryResolver) GetEvents(ctx context.Context) ([]*model.Event, error) {
	return r.db.GetAllEvents(ctx)
}

// GetTracks is the resolver for the getTracks field.
func (r *queryResolver) GetTracks(ctx context.Context) ([]*model.Track, error) {
	return r.db.GetAllTracks(ctx)
}

// Track is the resolver for the track field.
func (r *queryResolver) Track(ctx context.Context, id int) (*model.Track, error) {
	return dataloader.For(ctx).GetTrack(ctx, id)
	// return dataloader.For(ctx).GetTrackNew(ctx, id)
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context, ids []int) ([]*model.Event, error) {
	// ret, _ := r.db.GetEvents(ctx, ids)
	ret, _ := dataloader.For(ctx).GetEvents(ctx, ids)
	return ret, nil
}

// Tracks is the resolver for the tracks field.
func (r *queryResolver) Tracks(ctx context.Context, ids []int) ([]*model.Track, error) {
	// return r.db.GetTracks(ctx, id)
	ret, _ := dataloader.For(ctx).GetTracks(ctx, ids)
	return ret, nil
}

// SearchDriver is the resolver for the searchDriver field.
func (r *queryResolver) SearchDriver(ctx context.Context, arg string) ([]*model.Driver, error) {
	dbResults := r.db.SearchDrivers(ctx, arg)
	return dbResults, nil
}

// SearchTeam is the resolver for the searchTeam field.
func (r *queryResolver) SearchTeam(ctx context.Context, arg string) ([]*model.Team, error) {
	dbResults := r.db.SearchTeams(ctx, arg)
	return dbResults, nil
}

// Drivers is the resolver for the drivers field.
func (r *teamResolver) Drivers(ctx context.Context, obj *model.Team) ([]*model.Driver, error) {
	dbResults, _ := dataloader.For(ctx).GetTeamDrivers(ctx, obj.Name)
	// log.Printf("dbResult: %v err: %v\n", dbResults, err)
	return dbResults, nil
}

// Events is the resolver for the events field.
func (r *teamResolver) Events(ctx context.Context, obj *model.Team) ([]*model.Event, error) {
	// eventIds := r.db.CollectEventIdsForTeam(ctx, obj.Name)
	eventIds := dataloader.For(ctx).GetEventIdsForTeam(ctx, obj.Name)
	tmp, _ := dataloader.For(ctx).GetEvents(ctx, eventIds)
	return tmp, nil
}

// Events is the resolver for the events field.
func (r *trackResolver) Events(ctx context.Context, obj *model.Track) ([]*model.Event, error) {
	tmp := dataloader.For(ctx).GetEventsForTrack(ctx, obj.ID)
	return tmp, nil
}

// Driver returns generated.DriverResolver implementation.
func (r *Resolver) Driver() generated.DriverResolver { return &driverResolver{r} }

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Team returns generated.TeamResolver implementation.
func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

// Track returns generated.TrackResolver implementation.
func (r *Resolver) Track() generated.TrackResolver { return &trackResolver{r} }

type driverResolver struct{ *Resolver }
type eventResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }
type trackResolver struct{ *Resolver }
