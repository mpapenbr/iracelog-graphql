package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/mpapenbr/iracelog-graphql/graph/generated"
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	"github.com/mpapenbr/iracelog-graphql/internal/tracks"
)

// Track is the resolver for the track field.
func (r *eventResolver) Track(ctx context.Context, obj *model.Event) (*model.Track, error) {
	track := tracks.GetById(obj.TrackId)
	result := &model.Track{ID: fmt.Sprintf("%d", track.ID), Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length}
	return result, nil
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context) ([]*model.Event, error) {
	var result []*model.Event

	events := events.GetALl()
	for _, event := range events {
		result = append(result, &model.Event{ID: fmt.Sprintf("%d", event.ID), Name: event.Name, Key: event.Key, TrackId: int64(event.Info.TrackId)})
	}
	return result, nil
}

// Tracks is the resolver for the tracks field.
func (r *queryResolver) Tracks(ctx context.Context) ([]*model.Track, error) {
	var result []*model.Track

	tracks := tracks.GetALl()
	for _, track := range tracks {
		result = append(result, &model.Track{ID: fmt.Sprintf("%d", track.ID), Name: track.Data.Name, ShortName: track.Data.ShortName, Length: track.Data.Length})
	}
	return result, nil
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type eventResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
