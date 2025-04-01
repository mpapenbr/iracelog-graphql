package car

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
	"github.com/stephenafamo/scan"

	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
	"github.com/mpapenbr/iracelog-graphql/internal/utils"
)

func GetEventCars(exec bob.Executor, eventIDs []int) (map[int][]*models.CCar, error) {
	myIDs := utils.IntSliceToInt32Slice(eventIDs)

	query := models.CCars.Query(
		// Note: we use any(myIDs) instead of In(myIDs...) here
		// in Postgres the IN operator is limited to 32k elements
		// even if we probably never reach that limit, it's better to be safe
		// otherwise we could use models.SelectWhere.CCars.EventID.In(myIDs...),
		// bonus: for IN we have to check for empty ids, with any we don't
		// bonus: we learn how to code this with bob ;)
		sm.Where(models.CCarColumns.EventID.EQ(psql.F("ANY", expr.Arg(myIDs)))),
		sm.OrderBy(models.CCarColumns.Name),
	)

	res, err := query.All(context.Background(), exec)
	if err != nil {
		return nil, err
	}

	ret := map[int][]*models.CCar{}
	for i := range res {
		val, ok := ret[int(res[i].EventID)]
		if !ok {
			val = []*models.CCar{}
		}
		val = append(val, res[i])
		ret[int(res[i].EventID)] = val
	}

	return ret, nil
}

// see here for using bob with user defined structs
// https://github.com/ParkWithEase/parkeasy/blob/main/backend/internal/pkg/repositories/booking/postgres.go
//
//nolint:whitespace,lll // editor/linter issue
func GetEventEntryCars(
	exec bob.Executor,
	eventEntryIDs []int,
) (map[int]*models.CCar, error) {
	myIDs := utils.IntSliceToInt32Slice(eventEntryIDs)
	type myStruct struct {
		models.CCar
		EntryID int32 `db:"e_id"`
	}

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(models.CCars.Columns()),
		sm.Columns(models.CCarEntryColumns.ID.As("e_id")),
	}
	whereMods := []mods.Where[*dialect.SelectQuery]{
		sm.Where(models.CCarEntryColumns.ID.EQ(psql.F("ANY", expr.Arg(myIDs)))),
	}

	smods = append(smods,
		sm.From(models.TableNames.CCars),
		sm.InnerJoin(models.TableNames.CCarEntries).
			On(models.CCarEntryColumns.CCarID.EQ(models.CCarColumns.ID)),
		psql.WhereAnd(whereMods...),
	)

	query := psql.Select(smods...)

	res, err := bob.All(context.Background(), exec, query, scan.StructMapper[myStruct]())
	if err != nil {
		return nil, err
	}

	ret := map[int]*models.CCar{}
	for i := range res {
		ret[int(res[i].EntryID)] = &res[i].CCar
	}
	return ret, nil
}
