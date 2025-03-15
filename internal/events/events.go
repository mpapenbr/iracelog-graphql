package events

import (
	"context"
	"time"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/mods"

	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
)

type DbEvent struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Key                  string `json:"key"`
	Description          string
	EventTime            time.Time `json:"eventTime"`
	RaceloggerVersion    string    `json:"raceloggerVersion"`
	TeamRacing           bool      `json:"teamRacing"`
	MultiClass           bool      `json:"multiClass"`
	NumCarTypes          int       `json:"numCarTypes"`
	NumCarClasses        int       `json:"numCarClasses"`
	IrSessionId          int       `json:"irSessionId"`
	TrackId              int       `json:"trackId"`
	PitSpeed             float64   `json:"pitSpeed"`
	ReplayMinTimestamp   time.Time `json:"replayMinTimestamp"`
	ReplayMinSessionTime float64   `json:"replayMinSessionTime"`
	ReplayMaxSessionTime float64   `json:"replayMaxSessionTime"`
	Sessions             []Session `json:"sessions"`
}

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
func GetALl(exec bob.Executor, pageable internal.DbPageable) (
	models.EventSlice, error,
) {
	query := models.Events.Query()

	query.Apply(createPageableMods(pageable)...)

	ret, err := query.All(context.Background(), exec)
	return ret, err
}

//nolint:whitespace // editor/linter issue
func GetByIds(exec bob.Executor, ids []int) (
	models.EventSlice, error,
) {
	myIds := make([]int32, len(ids))
	for i, v := range ids {
		myIds[i] = int32(v)
	}
	if len(ids) == 0 {
		return models.EventSlice{}, nil
	}
	ret, err := models.Events.Query(
		models.SelectWhere.Events.ID.In(myIds...),
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
func GetEventsByTrackIds(
	exec bob.Executor,
	trackIds []int,
	pageable internal.DbPageable,
) (map[int][]*models.Event, error) {
	myIds := make([]int32, len(trackIds))
	for i, v := range trackIds {
		myIds[i] = int32(v)
	}
	if len(myIds) == 0 {
		return map[int][]*models.Event{}, nil
	}
	query := models.Events.Query(
		models.SelectWhere.Events.TrackID.In(myIds...),
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
	searchArg string,
	pageable internal.DbPageable,
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

	query := models.Events.Query(partMain)
	query.Apply(createPageableMods(pageable)...)
	res, err := query.All(context.Background(), exec)
	return res, err
}

//nolint:lll,funlen,whitespace // sql readability
func AdvancedEventSearch(
	exec bob.Executor,
	search *EventSearchKeys,
	pageable internal.DbPageable,
) (models.EventSlice, error) {
	query := models.Events.Query()
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
	exec = bob.Debug(exec)

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

func createPageableMods(pageable internal.DbPageable) []bob.Mod[*dialect.SelectQuery] {
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
