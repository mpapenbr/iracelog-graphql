package tracks

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
)

//nolint:tagliatelle // json is that way
type Sector struct {
	Num      int     `json:"num"`
	StartPct float64 `json:"start_pct"`
}

//nolint:whitespace // editor/linter issue
func GetAll(exec bob.Executor, pageable internal.DBPageable) (
	models.TrackSlice, error,
) {
	query := models.Tracks.Query()

	mods := make([]bob.Mod[*dialect.SelectQuery], 0)
	if pageable.Limit != nil {
		mods = append(mods, sm.Limit(*pageable.Limit))
	}
	if pageable.Offset != nil {
		mods = append(mods, sm.Offset(*pageable.Offset))
	}
	if pageable.Sort != nil {
		//nolint:gocritic // by design
		for _, s := range pageable.Sort.Expressions {
			mods = append(mods, sm.OrderBy(s))
		}
	}
	query.Apply(mods...)

	ret, err := query.All(context.Background(), exec)
	return ret, err
}

func GetByIDs(exec bob.Executor, ids []int) (models.TrackSlice, error) {
	myIDs := make([]int32, len(ids))
	for i, v := range ids {
		myIDs[i] = int32(v)
	}
	if len(ids) == 0 {
		return models.TrackSlice{}, nil
	}
	ret, err := models.Tracks.Query(
		models.SelectWhere.Tracks.ID.In(myIDs...),
	).All(context.Background(), exec)

	return ret, err
}
