package events

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/mods"

	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
)

//nolint:tagliatelle // json is that way
type Session struct {
	Num         int    `json:"num"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	SessionTime int    `json:"session_time"`
	Laps        int    `json:"laps"`
}

type EventSearchKeys struct {
	Name   string
	Car    string
	Track  string
	Driver string
	Team   string
}

//nolint:whitespace // editor/linter issue
func GetALl(exec bob.Executor, tenantID int, pageable internal.DBPageable) (
	models.EventSlice, error,
) {
	query := models.Events.Query(
		models.SelectWhere.Events.TenantID.EQ(int32(tenantID)),
	)

	query.Apply(createPageableMods(pageable)...)

	ret, err := query.All(context.Background(), exec)
	return ret, err
}

//nolint:whitespace // editor/linter issue
func GetByIDs(exec bob.Executor, tenantID int, ids []int) (
	models.EventSlice, error,
) {
	myIDs := make([]int32, len(ids))
	for i, v := range ids {
		myIDs[i] = int32(v)
	}
	if len(ids) == 0 {
		return models.EventSlice{}, nil
	}
	ret, err := models.Events.Query(
		models.SelectWhere.Events.ID.In(myIDs...),
		models.SelectWhere.Events.TenantID.EQ(int32(tenantID)),
	).All(context.Background(), exec)
	return ret, err
}

/*
note: currently only pageable.sort is processed.
Discussion: how should limit/offset be interpreted?
We can't put it on the query as this would limit/offset the overall data.
So we have to process it "manually" for each event, which yields the next question:
should offset apply only for those tracks having more than offset events?
consider a track with 2 and another with 10 events and a query with offset 5
*/
//nolint:whitespace // editor/linter issue
func GetEventsByTrackIDs(
	exec bob.Executor,
	tenantID int,
	trackIDs []int,
	pageable internal.DBPageable,
) (map[int][]*models.Event, error) {
	myIDs := make([]int32, len(trackIDs))
	for i, v := range trackIDs {
		myIDs[i] = int32(v)
	}
	if len(myIDs) == 0 {
		return map[int][]*models.Event{}, nil
	}
	query := models.Events.Query(
		models.SelectWhere.Events.TrackID.In(myIDs...),
		models.SelectWhere.Events.TenantID.EQ(int32(tenantID)),
	)
	query.Apply(createPageableMods(pageable)...)
	res, err := query.All(context.Background(), exec)
	if err != nil {
		return nil, err
	}

	ret := map[int][]*models.Event{}
	for i := range res {
		val, ok := ret[int(res[i].TrackID)]
		if !ok {
			val = []*models.Event{}
		}
		val = append(val, res[i])
		ret[int(res[i].TrackID)] = val
	}

	return ret, nil
}

//nolint:lll,whitespace,funlen // sql readability
func SimpleEventSearch(
	exec bob.Executor,
	tenantID int,
	searchArg string,
	pageable internal.DBPageable,
) (models.EventSlice, error) {
	// we keep the original query for reference
	_ = `
select id,name	
WHERE name ilike $1
OR    description ilike $1
OR    track_id in (select id from track where name ilike $1)
OR id in (select event_id from c_car where name ilike $1)
OR id in (select e.event_id from c_car_entry e join c_car_team t on t.c_car_entry_id=e.id and t.name ilike $1)
OR id in (select e.event_id from c_car_entry e join c_car_driver d on d.c_car_entry_id=e.id and d.name ilike $1)
		`

	partMain := psql.WhereOr(
		models.SelectWhere.Events.Name.ILike(sqlStringContains(searchArg)),
		models.SelectWhere.Events.Description.ILike(sqlStringContains(searchArg)),
		modSubQueryTrack(searchArg),
		modSubQueryCar(searchArg),
		modSubQueryTeam(searchArg),
		modSubQueryDriver(searchArg),
	)

	query := models.Events.Query(
		models.SelectWhere.Events.TenantID.EQ(int32(tenantID)),
		partMain,
	)
	query.Apply(createPageableMods(pageable)...)
	res, err := query.All(context.Background(), exec)
	return res, err
}

//nolint:lll,funlen,whitespace // sql readability
func AdvancedEventSearch(
	exec bob.Executor,
	tenantID int,
	search *EventSearchKeys,
	pageable internal.DBPageable,
) (models.EventSlice, error) {
	query := models.Events.Query(
		models.SelectWhere.Events.TenantID.EQ(int32(tenantID)),
	)
	query.Apply(createPageableMods(pageable)...)

	if search.Name != "" {
		w := models.SelectWhere.Events.Name.ILike(sqlStringContains(search.Name))
		query.Apply(w)
	}
	if search.Track != "" {
		query.Apply(modSubQueryTrack(sqlStringContains(search.Track)))
	}
	if search.Car != "" {
		query.Apply(modSubQueryCar(sqlStringContains(search.Car)))
	}
	if search.Driver != "" {
		query.Apply(modSubQueryDriver(sqlStringContains(search.Driver)))
	}
	if search.Team != "" {
		query.Apply(modSubQueryTeam(sqlStringContains(search.Team)))
	}

	res, err := query.All(context.Background(), exec)
	return res, err
}

func sqlStringContains(arg string) string {
	return "%" + arg + "%"
}

func modSubQueryTrack(searchArg string) mods.Where[*dialect.SelectQuery] {
	sub := psql.Select(
		sm.Columns(models.TrackColumns.ID),
		sm.From(models.TableNames.Tracks),
		models.SelectWhere.Tracks.Name.ILike(sqlStringContains(searchArg)),
	)
	w := sm.Where(models.EventColumns.TrackID.In(sub))
	return w
}

func modSubQueryCar(searchArg string) mods.Where[*dialect.SelectQuery] {
	sub := psql.Select(
		sm.Columns(models.CCarColumns.EventID),
		sm.From(models.TableNames.CCars),
		models.SelectWhere.CCars.Name.ILike(sqlStringContains(searchArg)),
	)
	w := sm.Where(models.EventColumns.ID.In(sub))
	return w
}

func modSubQueryDriver(searchArg string) mods.Where[*dialect.SelectQuery] {
	sub := psql.Select(
		sm.Columns(models.CCarEntryColumns.EventID),
		sm.From(models.TableNames.CCarEntries),
		sm.InnerJoin(models.TableNames.CCarDrivers).
			On(models.CCarEntryColumns.ID.EQ(models.CCarDriverColumns.CCarEntryID)),
		models.SelectWhere.CCarDrivers.Name.ILike(sqlStringContains(searchArg)))
	w := sm.Where(models.EventColumns.ID.In(sub))
	return w
}

func modSubQueryTeam(searchArg string) mods.Where[*dialect.SelectQuery] {
	sub := psql.Select(
		sm.Columns(models.CCarEntryColumns.EventID),
		sm.From(models.TableNames.CCarEntries),
		sm.InnerJoin(models.TableNames.CCarTeams).
			On(models.CCarEntryColumns.ID.EQ(models.CCarTeamColumns.CCarEntryID)),
		models.SelectWhere.CCarTeams.Name.ILike(sqlStringContains(searchArg)),
	)
	w := sm.Where(models.EventColumns.ID.In(sub))
	return w
}

func createPageableMods(pageable internal.DBPageable) []bob.Mod[*dialect.SelectQuery] {
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
	return mods
}
