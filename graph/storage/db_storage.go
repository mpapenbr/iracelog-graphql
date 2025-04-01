//nolint:dupl // can't be avoided
package storage

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"

	"github.com/mpapenbr/iracelog-graphql/graph/model"
	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
	"github.com/mpapenbr/iracelog-graphql/internal/tenant"
)

// for implementation of the storage interface see db_storage_xxx.go
// depending on the type of data to be returned
type (
	DBStorage struct {
		// Storage
		pool     *pgxpool.Pool
		executor bob.Executor
	}

	DBStorageOption func(*DBStorage)
)

func NewDBStorage(pool *pgxpool.Pool, opts ...DBStorageOption) Storage {
	db := stdlib.OpenDBFromPool(pool)

	ret := &DBStorage{
		pool:     pool,
		executor: bob.New(db),
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

// events
//
//nolint:whitespace // editor/linter issue
func (db *DBStorage) SimpleSearchEvents(
	ctx context.Context,
	arg string,
	limit *int,
	offset *int,
	sort []*model.EventSortArg,
) ([]*model.Event, error) {
	tp := tenant.GetFromContext(ctx)
	if tp == nil {
		return nil, tenant.ErrNoTenantInContext
	}
	tenantID, _ := tp()
	var result []*model.Event
	dbEventSortArg := convertEventSortArgs(sort)
	events, err := events.SimpleEventSearch(
		db.executor,
		tenantID,
		arg,
		internal.DBPageable{Limit: limit, Offset: offset, Sort: dbEventSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, dbEvents := range events {
			// this would cause assigning the last loop content to all result entries

			result = append(result, convertDBEventToModel(dbEvents))
		}
	}
	return result, err
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) AdvancedSearchEvents(
	ctx context.Context,
	arg *events.EventSearchKeys,
	limit *int,
	offset *int,
	sort []*model.EventSortArg,
) ([]*model.Event, error) {
	var result []*model.Event
	tp := tenant.GetFromContext(ctx)
	if tp == nil {
		return nil, tenant.ErrNoTenantInContext
	}
	tenantID, _ := tp()
	dbEventSortArg := convertEventSortArgs(sort)
	events, err := events.AdvancedEventSearch(
		db.executor,
		tenantID,
		arg,
		internal.DBPageable{Limit: limit, Offset: offset, Sort: dbEventSortArg})
	if err == nil {
		// convert the internal database Track to the GraphQL-Track
		for _, dbEvents := range events {
			// this would cause assigning the last loop content to all result entries

			result = append(result, convertDBEventToModel(dbEvents))
		}
	}
	return result, err
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectAnalysisData(
	ctx context.Context,
	eventIDs dataloader.Keys,
) map[string]analysis.DBAnalysis {
	res, _ := analysis.GetAnalysisForEvents(db.pool, IntKeysToSlice(eventIDs))
	ret := map[string]analysis.DBAnalysis{}
	//nolint:gocritic // by design
	for _, a := range res {
		ret[IntKey(a.EventID).String()] = a
	}
	return ret
}

func (db *DBStorage) SearchDrivers(ctx context.Context, arg string) []*model.Driver {
	res, _ := analysis.SearchDrivers(db.pool, arg)
	ret := make([]*model.Driver, len(res))
	//nolint:gocritic // by design
	for i, d := range res {
		teams := make([]*model.Team, len(d.Teams))
		for j, d := range d.Teams {
			teams[j] = &model.Team{Name: d}
		}
		ret[i] = &model.Driver{
			Name:     d.Name,
			Teams:    teams,
			CarNum:   d.CarNum,
			CarClass: d.CarClass,
		}
	}
	return ret
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectDriversInTeams(
	ctx context.Context,
	teams dataloader.Keys,
) map[string][]*model.Driver {
	res, _ := analysis.SearchDriversInTeams(db.pool, teams.Keys())
	ret := map[string][]*model.Driver{}
	for k, v := range res {
		drivers := make([]*model.Driver, len(v))
		//nolint:gocritic // by design
		for i, d := range v {
			drivers[i] = &model.Driver{
				Name:     d.Name,
				CarNum:   d.CarNum,
				CarClass: d.CarClass,
			}
		}
		ret[k] = drivers
	}
	return ret
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectTeamsForDrivers(
	ctx context.Context,
	drivers dataloader.Keys,
) map[string][]*model.Team {
	res, _ := analysis.SearchTeamsForDrivers(db.pool, drivers.Keys())
	ret := map[string][]*model.Team{}
	//nolint:gocritic // by design
	for k, v := range res {
		teams := make([]*model.Team, len(v))
		for i, d := range v {
			teams[i] = &model.Team{Name: d.Name, CarNum: d.CarNum, CarClass: d.CarClass}
		}
		ret[k] = teams
	}
	return ret
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectEventIDsForTeams(
	ctx context.Context,
	teams dataloader.Keys,
) map[string][]int {
	ret, _ := analysis.CollectEventIDsForTeams(db.pool, teams.Keys())
	return ret
}

//nolint:whitespace // editor/linter issue
func (db *DBStorage) CollectEventIDsForDrivers(
	ctx context.Context,
	drivers dataloader.Keys,
) map[string][]int {
	ret, _ := analysis.CollectEventIDsForDrivers(db.pool, drivers.Keys())
	return ret
}

func (db *DBStorage) SearchTeams(ctx context.Context, arg string) []*model.Team {
	res, _ := analysis.SearchTeams(db.pool, arg)
	ret := make([]*model.Team, len(res))
	//nolint:gocritic // by design
	for i, d := range res {
		drivers := make([]*model.Driver, len(d.Drivers))
		for j, d := range d.Drivers {
			drivers[j] = &model.Driver{Name: d}
		}
		ret[i] = &model.Team{
			Name:     d.Name,
			Drivers:  drivers,
			CarNum:   d.CarNum,
			CarClass: d.CarClass,
		}
	}
	return ret
}
