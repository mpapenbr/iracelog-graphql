package storage

import (
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
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

func convertDBEventToModel(dbEvent *models.Event) *model.Event {
	return &model.Event{
		ID:                int(dbEvent.ID),
		Name:              dbEvent.Name,
		Description:       dbEvent.Description.V,
		Key:               dbEvent.EventKey,
		TrackID:           int(dbEvent.TrackID),
		RecordDate:        dbEvent.EventTime,
		RaceloggerVersion: dbEvent.RaceloggerVersion,
		TeamRacing:        dbEvent.TeamRacing,
		MultiClass:        dbEvent.MultiClass,
		IracingSessionID:  int(dbEvent.IrSubSessionID),
		NumCarClasses:     int(dbEvent.NumCarClasses),
		NumCarTypes:       int(dbEvent.NumCarTypes),
	}
}

//nolint:errcheck // by design
func convertDBTrackToModel(dbTrack *models.Track) *model.Track {
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

func convertDBCarToModel(d *models.CCar) *model.Car {
	return &model.Car{
		ID:            int(d.ID),
		Name:          d.Name,
		NameShort:     d.NameShort,
		CarID:         int(d.CarID),
		FuelPct:       d.FuelPCT.InexactFloat64(),
		PowerAdjust:   d.PowerAdjust.InexactFloat64(),
		WeightPenalty: d.WeightPenalty.InexactFloat64(),
		DryTireSets:   int(d.DryTireSets),
	}
}

func convertDBEventEntryToModel(d *models.CCarEntry) *model.EventEntry {
	return &model.EventEntry{
		ID:        int(d.ID),
		CarNum:    &d.CarNumber,
		CarNumRaw: toIntPtr(d.CarNumberRaw),
	}
}

func convertDBCarTeamToModel(d *models.CCarTeam) *model.EventTeam {
	return &model.EventTeam{
		ID:     int(d.ID),
		Name:   d.Name,
		TeamID: int(d.TeamID),
	}
}

func convertDBCarDriverToModel(d *models.CCarDriver) *model.EventDriver {
	return &model.EventDriver{
		ID:              int(d.ID),
		Name:            d.Name,
		DriverID:        int(d.DriverID),
		Initials:        &d.Initials,
		AbbrevName:      &d.AbbrevName,
		IRating:         toIntPtr(d.Irating),
		LicenseLevel:    toIntPtr(d.LicLevel),
		LicenseSubLevel: toIntPtr(d.LicSubLevel),
		LicenseString:   &d.LicString,
	}
}

func toIntPtr(i int32) *int {
	ret := int(i)
	return &ret
}
