package storage

import (
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
