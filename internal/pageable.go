package internal

import (
	"fmt"
	"strings"
)

type DbPageable struct {
	Limit  *int
	Offset *int
	Sort   []DbSortArg
}

type DbSortArg struct {
	Column string
	Order  string
}

func convertSortArg(args []DbSortArg) string {
	var ret []string
	for _, val := range args {
		ret = append(ret, fmt.Sprintf("%s %s", val.Column, val.Order))
	}
	return strings.Join(ret, ",")
}

func HandlePageableArgs(query string, pageable DbPageable) string {
	ret := query
	if len(pageable.Sort) > 0 {
		ret = fmt.Sprintf("%s order by %s", ret, convertSortArg(pageable.Sort))
	}

	if pageable.Offset != nil {
		ret = fmt.Sprintf("%s offset %d", ret, *pageable.Offset)
	}
	if pageable.Limit != nil && *pageable.Limit > 0 {
		ret = fmt.Sprintf("%s limit %d", ret, *pageable.Limit)
	}
	return ret
}
