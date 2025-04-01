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
func (db *DBStorage) GetAllEvents(
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
	tenantID, _ := tp()
	dbEventSortArg := convertEventSortArgs(sort)
	events, err := events.GetALl(db.executor, tenantID, internal.DBPageable{
		Limit:  limit,
		Offset: offset,
		Sort:   dbEventSortArg,
	})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, dbEvents := range events {
			// this would cause assigning the last loop content to all result entries

			result = append(result, convertDBEventToModel(dbEvents))
		}
	}
	return result, err
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) GetEventsByKeys(
	ctx context.Context,
	ids dataloader.Keys,
) map[string]*model.Event {
	tp := tenant.GetFromContext(ctx)
	if tp == nil {
		return nil
	}
	tenantID, _ := tp()
	intIDs := IntKeysToSlice(ids)
	result := map[string]*model.Event{}
	events, _ := events.GetByIDs(db.executor, tenantID, intIDs)
	// convert the internal database Track to the GraphQL-Track
	for _, dbEvents := range events {
		// this would cause assigning the last loop content to all result entries
		result[IntKey(dbEvents.ID).String()] = convertDBEventToModel(dbEvents)
	}

	return result
}

// Note: we use (temporary) a string as key (to reuse existing batcher mechanics)
//
//nolint:whitespace // editor/linter issue
func (db *DBStorage) GetEventsForTrackIDsKeys(
	ctx context.Context,
	trackIDs dataloader.Keys,
) map[string][]*model.Event {
	tp := tenant.GetFromContext(ctx)
	if tp == nil {
		return nil
	}
	tenantID, _ := tp()
	result := map[string][]*model.Event{}

	intTrackIDs := make([]int, len(trackIDs))
	for i, id := range trackIDs {
		//nolint:errcheck // we know that the conversion works
		intTrackIDs[i] = id.Raw().(int)
	}
	byTrackID, err := events.GetEventsByTrackIDs(
		db.executor,
		tenantID,
		intTrackIDs,
		internal.DBPageable{Sort: convertEventSortArgs([]*model.EventSortArg{})})
	if err == nil {
		// convert the internal database Event to the GraphQL-Event
		for k, event := range byTrackID {
			convertedEvents := make([]*model.Event, len(event))
			for i, dbEvent := range event {
				convertedEvents[i] = convertDBEventToModel(dbEvent)
			}
			result[fmt.Sprintf("%d", k)] = convertedEvents
		}
	}
	return result
}
