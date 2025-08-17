package entry

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/expr"

	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
	"github.com/mpapenbr/iracelog-graphql/internal/utils"
)

//nolint:whitespace // editor/linter issue
func GetEventEntriesByEventID(
	exec bob.Executor,
	eventIDs []int,
) (map[int][]*models.CCarEntry, error) {
	myIDs := utils.IntSliceToInt32Slice(eventIDs)

	query := models.CCarEntries.Query(
		// see also notes in interanl/car/car/car.go
		sm.Where(models.CCarEntries.Columns.EventID.EQ(psql.F("ANY", expr.Arg(myIDs)))),
		sm.OrderBy(models.CCarEntries.Columns.CarNumberRaw),
	)

	res, err := query.All(context.Background(), exec)
	if err != nil {
		return nil, err
	}

	ret := map[int][]*models.CCarEntry{}
	for i := range res {
		val, ok := ret[int(res[i].EventID)]
		if !ok {
			val = []*models.CCarEntry{}
		}
		val = append(val, res[i])
		ret[int(res[i].EventID)] = val
	}

	return ret, nil
}

//nolint:whitespace // editor/linter issue
func GetEventEntriesByIDs(
	exec bob.Executor,
	ids []int) (
	map[int]*models.CCarEntry, error,
) {
	myIDs := utils.IntSliceToInt32Slice(ids)

	query := models.CCarEntries.Query(
		// see also notes in interanl/car/car/car.go
		sm.Where(models.CCarEntries.Columns.ID.EQ(psql.F("ANY", expr.Arg(myIDs)))),
		sm.OrderBy(models.CCarEntries.Columns.CarNumberRaw),
	)

	res, err := query.All(context.Background(), exec)
	if err != nil {
		return nil, err
	}

	ret := map[int]*models.CCarEntry{}
	for i := range res {
		ret[int(res[i].EventID)] = res[i]
	}

	return ret, nil
}
