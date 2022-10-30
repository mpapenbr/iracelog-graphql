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

func ConvertSortArg(args []DbSortArg) string {
	var ret []string
	for _, val := range args {
		ret = append(ret, fmt.Sprintf("%s %s", val.Column, val.Order))
	}
	return strings.Join(ret, ",")
}
