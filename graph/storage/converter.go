package storage

import (
	"time"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
)

// converts model arguments to db arguments

func convertEventSortArgs(modelArgs []*model.EventSortArg) []events.DbEventSortArg {
	if modelArgs == nil || len(modelArgs) == 0 {
		ret := []events.DbEventSortArg{
			// {Column: "data->'info'->'trackDisplayName'", Order: "asc"},
			{Column: "id", Order: "desc"},
		}
		return ret
	}
	var ret []events.DbEventSortArg
	for _, arg := range modelArgs {
		var item events.DbEventSortArg
		switch arg.Field {
		case model.EventSortFieldName:
			item.Column = "name"
		case model.EventSortFieldRecordDate:
			item.Column = "record_stamp"
		case model.EventSortFieldTrack:
			item.Column = "data->'info'->'trackDisplayName'"
		}
		if arg.Order != nil && *arg.Order == model.SortOrderDesc {
			item.Order = "desc"
		} else {
			item.Order = "asc"
		}

		ret = append(ret, item)
	}

	return ret

}

func convertDbEventToModel(dbEvent events.DbEvent) *model.Event {
	eventTime, _ := time.Parse("2006-01-02T15:04:05", dbEvent.Info.EventTime)

	return &model.Event{
		ID:                dbEvent.ID,
		Name:              dbEvent.Name,
		Key:               dbEvent.Key,
		TrackId:           dbEvent.Info.TrackId,
		RecordDate:        dbEvent.RecordStamp,
		EventDate:         eventTime,
		RaceloggerVersion: dbEvent.Info.RaceloggerVersion,
		TeamRacing:        dbEvent.Info.TeamRacing == 1,
		MultiClass:        dbEvent.Info.MultiClass,
		IracingSessionId:  dbEvent.Info.IrSessionId,
		NumCarClasses:     dbEvent.Info.NumCarClasses,
		NumCarTypes:       dbEvent.Info.NumCarTypes,
		Track:             &model.Track{},
		DbEvent:           &dbEvent,
	}
}
