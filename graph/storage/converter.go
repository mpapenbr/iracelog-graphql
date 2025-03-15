package storage

import (
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal/car/car"
	"github.com/mpapenbr/iracelog-graphql/internal/car/driver"
	"github.com/mpapenbr/iracelog-graphql/internal/car/entry"
	"github.com/mpapenbr/iracelog-graphql/internal/car/team"
	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

func convertEventSortArgs(modelArgs []*model.EventSortArg) *clause.OrderBy {
	ret := clause.OrderBy{}
	if len(modelArgs) == 0 {
		ret.AppendOrder(clause.OrderDef{
			Expression: models.EventColumns.EventTime,
			Direction:  DESC,
		})

		return &ret
	}
	for _, arg := range modelArgs {
		var item clause.OrderDef
		//nolint:exhaustive // by design
		switch arg.Field {
		case model.EventSortFieldName:
			item.Expression = models.EventColumns.Name
		case model.EventSortFieldRecordDate:
			item.Expression = models.EventColumns.EventTime
		}
		if arg.Order != nil && *arg.Order == model.SortOrderDesc {
			item.Direction = DESC
		} else {
			item.Direction = ASC
		}
		ret.AppendOrder(item)
	}
	return &ret
}

func convertTrackSortArgs(modelArgs []*model.TrackSortArg) *clause.OrderBy {
	ret := clause.OrderBy{}
	if len(modelArgs) == 0 {
		ret.AppendOrder(clause.OrderDef{
			Expression: models.TrackColumns.Name,
			Direction:  ASC,
		})

		return &ret
	}
	for _, arg := range modelArgs {
		var item clause.OrderDef
		switch arg.Field {
		case model.TrackSortFieldName:
			item.Expression = models.TrackColumns.Name
		case model.TrackSortFieldShortName:
			item.Expression = models.TrackColumns.ShortName
		case model.TrackSortFieldID:
			item.Expression = models.TrackColumns.ID
		case model.TrackSortFieldLength:
			item.Expression = models.TrackColumns.TrackLength
		case model.TrackSortFieldPitlaneLength:
			item.Expression = models.TrackColumns.PitLaneLength
		case model.TrackSortFieldNumSectors:
			item.Expression = dialect.NewExpression(
				psql.F("jsonb_array_length", models.TrackColumns.Sectors))
		}
		if arg.Order != nil && *arg.Order == model.SortOrderDesc {
			item.Direction = DESC
		} else {
			item.Direction = ASC
		}
		ret.AppendOrder(item)
	}
	return &ret
}

func convertDbEventToModel(dbEvent *models.Event) *model.Event {
	return &model.Event{
		ID:                int(dbEvent.ID),
		Name:              dbEvent.Name,
		Description:       dbEvent.Description.GetOr(""),
		Key:               dbEvent.EventKey,
		TrackId:           int(dbEvent.TrackID),
		RecordDate:        dbEvent.EventTime,
		RaceloggerVersion: dbEvent.RaceloggerVersion,
		TeamRacing:        dbEvent.TeamRacing,
		MultiClass:        dbEvent.MultiClass,
		IracingSessionId:  int(dbEvent.IrSubSessionID),
		NumCarClasses:     int(dbEvent.NumCarClasses),
		NumCarTypes:       int(dbEvent.NumCarTypes),
	}
}

//nolint:errcheck // by design
func convertDbTrackToModel(dbTrack *models.Track) *model.Track {
	return &model.Track{
		ID:            int(dbTrack.ID),
		Name:          dbTrack.Name,
		ShortName:     dbTrack.ShortName,
		ConfigName:    dbTrack.Config,
		Length:        dbTrack.TrackLength.Abs().InexactFloat64(),
		NumSectors:    len(dbTrack.Sectors),
		PitLaneLength: dbTrack.PitLaneLength.Abs().InexactFloat64(),
		PitSpeed:      dbTrack.PitSpeed.Abs().InexactFloat64(),
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
