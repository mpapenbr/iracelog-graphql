package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/mpapenbr/iracelog-graphql/graph/dataloader"
	"github.com/mpapenbr/iracelog-graphql/graph/generated"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
)

// Track is the resolver for the track field.
func (r *eventResolver) Track(ctx context.Context, obj *model.Event) (*model.Track, error) {
	// track := tracks.GetById(obj.TrackId)
	// result := &model.Track{ID: track.ID, Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length}
	// fmt.Printf("eventResolver.Track, event=%v, track=%v\n", obj.ID, obj.TrackId)
	return dataloader.For(ctx).GetTrack(ctx, obj.TrackId)
}

// Teams is the resolver for the teams field.
func (r *eventResolver) Teams(ctx context.Context, obj *model.Event) ([]*model.Team, error) {
	return r.db.GetTeamsForEvent(ctx, obj), nil
}

// Drivers is the resolver for the drivers field.
func (r *eventResolver) Drivers(ctx context.Context, obj *model.Event) ([]*model.Driver, error) {
	panic(fmt.Errorf("not implemented"))
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
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context, ids []int) ([]*model.Event, error) {
	ret, _ := r.db.GetEvents(ctx, ids)
	return ret, nil
}

// Tracks is the resolver for the tracks field.
func (r *queryResolver) Tracks(ctx context.Context, ids []int) ([]*model.Track, error) {
	// return r.db.GetTracks(ctx, id)
	ret, _ := dataloader.For(ctx).GetTracks(ctx, ids)
	return ret, nil
}

// Events is the resolver for the events field.
func (r *trackResolver) Events(ctx context.Context, obj *model.Track) ([]*model.Event, error) {
	eventIds, _ := r.db.GetEventIdsForTrackId(ctx, obj.ID)

	tmp, err := dataloader.For(ctx).GetEvents(ctx, eventIds)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, errors.New("Dings")
	}
	return tmp, nil
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Track returns generated.TrackResolver implementation.
func (r *Resolver) Track() generated.TrackResolver { return &trackResolver{r} }

type eventResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type trackResolver struct{ *Resolver }
