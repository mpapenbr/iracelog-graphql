package storage

import (
	"time"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	"github.com/mpapenbr/iracelog-graphql/internal/tracks"
)

// converts model arguments to db arguments

func convertEventSortArgs(modelArgs []*model.EventSortArg) []internal.DbSortArg {
	if len(modelArgs) == 0 {
		ret := []internal.DbSortArg{
			{Column: "record_stamp", Order: "desc"},
		}
		return ret
	}
	var ret []internal.DbSortArg
	for _, arg := range modelArgs {
		var item internal.DbSortArg
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

func convertTrackSortArgs(modelArgs []*model.TrackSortArg) []internal.DbSortArg {
	if len(modelArgs) == 0 {
		ret := []internal.DbSortArg{
			{Column: "data->>'trackDisplayName'", Order: "asc"},
		}
		return ret
	}
	var ret []internal.DbSortArg
	for _, arg := range modelArgs {
		var item internal.DbSortArg
		switch arg.Field {

		case model.TrackSortFieldName:
			item.Column = "data->>'trackDisplayName'"
		case model.TrackSortFieldShortName:
			item.Column = "data->>'trackDisplayShortName'"
		case model.TrackSortFieldID:
			item.Column = "id"
		case model.TrackSortFieldLength:
			item.Column = "data->'trackLength'"
		case model.TrackSortFieldPitlaneLength:
			item.Column = "data->'pit'->'laneLength'" // TODO
		case model.TrackSortFieldNumSectors:
			item.Column = "jsonb_array_length(data->'sectors')"
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
		Description:       dbEvent.Description,
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

func convertDbTrackToModel(dbTrack tracks.DbTrack) *model.Track {

	return &model.Track{
		ID:            dbTrack.ID,
		Name:          dbTrack.Data.Name,
		ShortName:     dbTrack.Data.ShortName,
		ConfigName:    dbTrack.Data.Config,
		Length:        dbTrack.Data.Length,
		NumSectors:    len(dbTrack.Data.Sectors),
		PitlaneLength: dbTrack.Data.Pit.LaneLength,
	}
}
