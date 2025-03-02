package storage

import (
	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/car/car"
	"github.com/mpapenbr/iracelog-graphql/internal/car/driver"
	"github.com/mpapenbr/iracelog-graphql/internal/car/entry"
	"github.com/mpapenbr/iracelog-graphql/internal/car/team"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	"github.com/mpapenbr/iracelog-graphql/internal/tracks"
)

// converts model arguments to db arguments

func convertEventSortArgs(modelArgs []*model.EventSortArg) []internal.DbSortArg {
	if len(modelArgs) == 0 {
		ret := []internal.DbSortArg{
			{Column: "event_time", Order: "desc"},
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
			item.Column = "event_time"
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
			{Column: "name'", Order: "asc"},
		}
		return ret
	}
	var ret []internal.DbSortArg
	for _, arg := range modelArgs {
		var item internal.DbSortArg
		switch arg.Field {
		case model.TrackSortFieldName:
			item.Column = "name'"
		case model.TrackSortFieldShortName:
			item.Column = "short_name'"
		case model.TrackSortFieldID:
			item.Column = "id"
		case model.TrackSortFieldLength:
			item.Column = "track_length"
		case model.TrackSortFieldPitlaneLength:
			item.Column = "pit_lane_length"
		case model.TrackSortFieldNumSectors:
			item.Column = "jsonb_array_length(sectors)"
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

func convertDbEventToModel(dbEvent *events.DbEvent) *model.Event {
	return &model.Event{
		ID:                dbEvent.ID,
		Name:              dbEvent.Name,
		Description:       dbEvent.Description,
		Key:               dbEvent.Key,
		TrackId:           dbEvent.TrackId,
		RecordDate:        dbEvent.EventTime,
		RaceloggerVersion: dbEvent.RaceloggerVersion,
		TeamRacing:        dbEvent.TeamRacing,
		MultiClass:        dbEvent.MultiClass,
		IracingSessionId:  dbEvent.IrSessionId,
		NumCarClasses:     dbEvent.NumCarClasses,
		NumCarTypes:       dbEvent.NumCarTypes,
		Track:             &model.Track{},
		DbEvent:           dbEvent,
	}
}

func convertDbTrackToModel(dbTrack *tracks.DbTrack) *model.Track {
	return &model.Track{
		ID:            dbTrack.ID,
		Name:          dbTrack.Name,
		ShortName:     dbTrack.ShortName,
		ConfigName:    dbTrack.Config,
		Length:        dbTrack.Length,
		NumSectors:    len(dbTrack.Sectors),
		PitLaneLength: dbTrack.PitLaneLength,
		PitSpeed:      dbTrack.PitSpeed,
	}
}

func convertDbCarToModel(d *car.DbCar) *model.Car {
	return &model.Car{
		ID:            d.ID,
		Name:          d.Name,
		NameShort:     d.NameShort,
		CarID:         d.CarId,
		FuelPct:       d.FuelPct,
		PowerAdjust:   d.PowerAdjust,
		WeightPenalty: d.WeightPenalty,
		DryTireSets:   d.DryTireSets,
	}
}

func convertDbEventEntryToModel(d *entry.DbCarEntry) *model.EventEntry {
	return &model.EventEntry{
		ID:        d.ID,
		CarNum:    &d.CarNum,
		CarNumRaw: &d.CarNumRaw,
	}
}

func convertDbCarTeamToModel(d *team.DbCarTeam) *model.EventTeam {
	return &model.EventTeam{
		ID:     d.ID,
		Name:   d.Name,
		TeamID: d.TeamId,
	}
}

func convertDbCarDriverToModel(d *driver.DbCarDriver) *model.EventDriver {
	return &model.EventDriver{
		ID:              d.ID,
		Name:            d.Name,
		DriverID:        d.DriverId,
		Initials:        &d.Initials,
		AbbrevName:      &d.AbbrevName,
		IRating:         &d.IRating,
		LicenseLevel:    &d.LicenseLevel,
		LicenseSubLevel: &d.LicenseSubLevel,
		LicenseString:   &d.LicenseString,
	}
}
