package team

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
func GetTeamsByEventEntry(
	exec bob.Executor,
	eventEntryIDs []int,
) (map[int]*models.CCarTeam, error) {
	myIds := utils.IntSliceToInt32Slice(eventEntryIDs)
	type myStruct struct {
		models.CCarTeam
		EntryId int32 `db:"e_id"`
	}

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(models.CCarTeams.Columns()),
		sm.Columns(models.CCarEntryColumns.ID.As("e_id")),
	}
	whereMods := []mods.Where[*dialect.SelectQuery]{
		sm.Where(models.CCarEntryColumns.ID.EQ(psql.F("ANY", expr.Arg(myIds)))),
	}

	smods = append(smods,
		sm.From(models.TableNames.CCarTeams),
		sm.InnerJoin(models.TableNames.CCarEntries).
			On(models.CCarEntryColumns.CCarID.EQ(models.CCarTeamColumns.ID)),
		psql.WhereAnd(whereMods...),
	)

	query := psql.Select(smods...)
	exec = bob.Debug(exec)
	res, err := bob.All(context.Background(), exec, query, scan.StructMapper[myStruct]())
	if err != nil {
		return nil, err
	}

	ret := map[int]*models.CCarTeam{}
	for i := range res {
		ret[int(res[i].EntryId)] = &res[i].CCarTeam
	}
	return ret, nil
}
