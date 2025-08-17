//nolint:dupl // ok until change to query builder
package driver

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

//nolint:whitespace // editor/linter issue
func GetEventDrivers(ctx context.Context, exec bob.Executor, eventIDs []int) (
	map[int][]*models.CCarDriver, error,
) {
	myIDs := utils.IntSliceToInt32Slice(eventIDs)
	type myStruct struct {
		models.CCarDriver
		EventID int32 `db:"event_id"`
	}

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(models.CCarDrivers.Columns.Names()),
		sm.Columns(models.CCarEntries.Columns.EventID.As("event_id")),
		sm.OrderBy(models.CCarDrivers.Columns.Name).Asc(),
	}
	whereMods := []mods.Where[*dialect.SelectQuery]{
		sm.Where(models.CCarEntries.Columns.EventID.EQ(psql.F("ANY", expr.Arg(myIDs)))),
	}

	smods = append(smods,
		sm.From(models.CCarDrivers.Name()),
		sm.InnerJoin(models.CCarEntries.Name()).
			On(models.CCarEntries.Columns.ID.EQ(models.CCarDrivers.Columns.CCarEntryID)),
		psql.WhereAnd(whereMods...),
	)

	query := psql.Select(smods...)

	res, err := bob.All(context.Background(), exec, query, scan.StructMapper[myStruct]())
	if err != nil {
		return nil, err
	}

	ret := map[int][]*models.CCarDriver{}
	for i := range res {
		val, ok := ret[int(res[i].EventID)]
		if !ok {
			val = []*models.CCarDriver{}
		}
		val = append(val, &res[i].CCarDriver)
		ret[int(res[i].EventID)] = val
	}
	return ret, nil
}

//nolint:whitespace // editor/linter issue
func GetDriversByEventEntry(ctx context.Context, exec bob.Executor, eventIDs []int) (
	map[int][]*models.CCarDriver, error,
) {
	myIDs := utils.IntSliceToInt32Slice(eventIDs)
	type myStruct struct {
		models.CCarDriver
		EntryID int32 `db:"e_id"`
	}

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(models.CCarDrivers.Columns.Names()),
		sm.Columns(models.CCarEntries.Columns.ID.As("e_id")),
		sm.OrderBy(models.CCarDrivers.Columns.Name).Asc(),
	}
	whereMods := []mods.Where[*dialect.SelectQuery]{
		sm.Where(models.CCarEntries.Columns.ID.EQ(psql.F("ANY", expr.Arg(myIDs)))),
	}

	smods = append(smods,
		sm.From(models.CCarDrivers.Name()),
		sm.InnerJoin(models.CCarEntries.Name()).
			On(models.CCarEntries.Columns.ID.EQ(models.CCarDrivers.Columns.CCarEntryID)),
		psql.WhereAnd(whereMods...),
	)

	query := psql.Select(smods...)

	res, err := bob.All(context.Background(), exec, query, scan.StructMapper[myStruct]())
	if err != nil {
		return nil, err
	}

	ret := map[int][]*models.CCarDriver{}
	for i := range res {
		val, ok := ret[int(res[i].EntryID)]
		if !ok {
			val = []*models.CCarDriver{}
		}
		val = append(val, &res[i].CCarDriver)
		ret[int(res[i].EntryID)] = val
	}
	return ret, nil
}

//nolint:whitespace,funlen // editor/linter issue
func GetDriversByTeam(
	ctx context.Context,
	exec bob.Executor,
	teamIDs []int,
) (
	map[int][]*models.CCarDriver, error,
) {
	myIDs := utils.IntSliceToInt32Slice(teamIDs)
	type myStruct struct {
		models.CCarDriver
		TeamID int32 `db:"t_id"`
	}

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(models.CCarDrivers.Columns.Names()),
		sm.Columns(models.CCarTeams.Columns.ID.As("t_id")),
		sm.OrderBy(models.CCarDrivers.Columns.Name).Asc(),
	}
	whereMods := []mods.Where[*dialect.SelectQuery]{
		sm.Where(models.CCarTeams.Columns.ID.EQ(psql.F("ANY", expr.Arg(myIDs)))),
	}

	smods = append(smods,
		sm.From(models.CCarDrivers.Name()),
		sm.InnerJoin(models.CCarEntries.Name()).
			On(models.CCarEntries.Columns.ID.EQ(models.CCarDrivers.Columns.CCarEntryID)),
		sm.InnerJoin(models.CCarTeams.Name()).
			On(models.CCarEntries.Columns.ID.EQ(models.CCarTeams.Columns.CCarEntryID)),
		psql.WhereAnd(whereMods...),
	)

	query := psql.Select(smods...)

	res, err := bob.All(context.Background(), exec, query, scan.StructMapper[myStruct]())
	if err != nil {
		return nil, err
	}

	ret := map[int][]*models.CCarDriver{}
	for i := range res {
		val, ok := ret[int(res[i].TeamID)]
		if !ok {
			val = []*models.CCarDriver{}
		}
		val = append(val, &res[i].CCarDriver)
		ret[int(res[i].TeamID)] = val
	}
	return ret, nil
}
