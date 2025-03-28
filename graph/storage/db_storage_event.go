package storage

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	"github.com/mpapenbr/iracelog-graphql/internal/tenant"
)

// contains implementations of storage interface that return a model.Event items
//
//nolint:whitespace // editor/linter issue
func (db *DbStorage) GetAllEvents(
	ctx context.Context,
	limit *int,
	offset *int,
	sort []*model.EventSortArg,
) ([]*model.Event, error) {
	var result []*model.Event
	tp := tenant.GetFromContext(ctx)
	if tp == nil {
		return nil, fmt.Errorf("tenant not found in context")
	}
	tenantId, _ := tp()
	dbEventSortArg := convertEventSortArgs(sort)
	events, err := events.GetALl(db.executor, tenantId, internal.DbPageable{
		Limit:  limit,
		Offset: offset,
		Sort:   dbEventSortArg,
	})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, dbEvents := range events {
			// this would cause assigning the last loop content to all result entries

			result = append(result, convertDbEventToModel(dbEvents))
		}
	}
	return result, err
}

//nolint:whitespace // editor/linter issue
func (db *DbStorage) GetEventsByKeys(
	ctx context.Context,
	ids dataloader.Keys,
) map[string]*model.Event {
	tp := tenant.GetFromContext(ctx)
	if tp == nil {
		return nil
	}
	tenantId, _ := tp()
	intIds := IntKeysToSlice(ids)
	result := map[string]*model.Event{}
	events, _ := events.GetByIds(db.executor, tenantId, intIds)
	// convert the internal database Track to the GraphQL-Track
	for _, dbEvents := range events {
		// this would cause assigning the last loop content to all result entries
		result[IntKey(dbEvents.ID).String()] = convertDbEventToModel(dbEvents)
	}

	return result
}

// Note: we use (temporary) a string as key (to reuse existing batcher mechanics)
//
//nolint:whitespace // editor/linter issue
func (db *DbStorage) GetEventsForTrackIdsKeys(
	ctx context.Context,
	trackIds dataloader.Keys,
) map[string][]*model.Event {
	tp := tenant.GetFromContext(ctx)
	if tp == nil {
		return nil
	}
	tenantId, _ := tp()
	result := map[string][]*model.Event{}

	intTrackIds := make([]int, len(trackIds))
	for i, id := range trackIds {
		//nolint:errcheck // we know that the conversion works
		intTrackIds[i] = id.Raw().(int)
	}
	byTrackId, err := events.GetEventsByTrackIds(
		db.executor,
		tenantId,
		intTrackIds,
		internal.DbPageable{Sort: convertEventSortArgs([]*model.EventSortArg{})})
	if err == nil {
		// convert the internal database Event to the GraphQL-Event
		for k, event := range byTrackId {
			convertedEvents := make([]*model.Event, len(event))
			for i, dbEvent := range event {
				convertedEvents[i] = convertDbEventToModel(dbEvent)
			}
			result[fmt.Sprintf("%d", k)] = convertedEvents
		}
	}
	return result
}
