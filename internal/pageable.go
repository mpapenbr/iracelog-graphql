package internal

import (
	"fmt"
	"strings"

	"github.com/stephenafamo/bob/clause"
)

type DBPageable struct {
	Limit   *int
	Offset  *int
	SortOld []DBSortArg
	Sort    *clause.OrderBy
}

type DBSortArg struct {
	Column string
	Order  string
}

func convertSortArg(args []DBSortArg) string {
	var ret []string
	//nolint:gocritic // by design
	for _, val := range args {
		ret = append(ret, fmt.Sprintf("%s %s", val.Column, val.Order))
	}
	return strings.Join(ret, ",")
}

func HandlePageableArgs(query string, pageable DBPageable) string {
	ret := query
	if len(pageable.SortOld) > 0 {
		ret = fmt.Sprintf("%s order by %s", ret, convertSortArg(pageable.SortOld))
	}

	if pageable.Offset != nil {
		ret = fmt.Sprintf("%s offset %d", ret, *pageable.Offset)
	}
	if pageable.Limit != nil && *pageable.Limit > 0 {
		ret = fmt.Sprintf("%s limit %d", ret, *pageable.Limit)
	}
	return ret
}
