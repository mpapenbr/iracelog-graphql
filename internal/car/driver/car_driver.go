//nolint:dupl // ok until change to query builder
package driver

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/mpapenbr/iracelog-graphql/internal"
)

type DbCarDriver struct {
	ID              int    `json:"id"`
	CarEntryId      int    `json:"carEntryId"`
	DriverId        int    `json:"driverId"`
	Name            string `json:"name"`
	Initials        string `json:"initials"`
	AbbrevName      string `json:"abbrevName"`
	IRating         int    `json:"iRating"`
	LicenseLevel    int    `json:"licenseLevel"`
	LicenseSubLevel int    `json:"licenseSubLevel"`
	LicenseString   string `json:"licenseString"`
}

//nolint:whitespace // editor/linter issue
func GetEventDrivers(ctx context.Context, q internal.Querier, eventIDs []int) (
	map[int][]*DbCarDriver, error,
) {
	rows, err := q.Query(ctx, `
	select d.id, d.c_car_entry_id, d.driver_id, d.name, d.initials, 
	  d.abbrev_name, d.irating, d.lic_level,d.lic_sub_level,
	  d.lic_string, e.event_id
	from c_car_driver d join c_car_entry e on d.c_car_entry_id = e.id
	where e.event_id = any($1) order by d.name asc`, eventIDs)
	if err != nil {
		log.Printf("error reading drivers: %v", err)
		return map[int][]*DbCarDriver{}, err
	}
	defer rows.Close()
	ret := map[int][]*DbCarDriver{}
	for rows.Next() {
		d := DbCarDriver{}
		var eventId int
		err := rows.Scan(&d.ID, &d.CarEntryId, &d.DriverId, &d.Name, &d.Initials,
			&d.AbbrevName, &d.IRating, &d.LicenseLevel, &d.LicenseSubLevel,
			&d.LicenseString, &eventId,
		)
		if err != nil {
			log.Printf("Error scaning driver: %v\n", err)
		}

		if _, ok := ret[eventId]; !ok {
			ret[eventId] = []*DbCarDriver{}
		}
		ret[eventId] = append(ret[eventId], &d)

	}
	return ret, nil
}

//nolint:whitespace // editor/linter issue
func GetDriversByEventEntry(
	ctx context.Context,
	q internal.Querier,
	eventEntryIDs []int) (
	map[int][]*DbCarDriver, error,
) {
	rows, err := q.Query(ctx, `
	select d.id, d.c_car_entry_id, d.driver_id, d.name, d.initials,
	 d.abbrev_name, d.irating, d.lic_level,d.lic_sub_level,
	 d.lic_string, e.id
	from c_car_driver d join c_car_entry e on d.c_car_entry_id = e.id
	where e.id = any($1) order by d.name asc`, eventEntryIDs)
	if err != nil {
		log.Printf("error reading drivers: %v", err)
		return map[int][]*DbCarDriver{}, err
	}
	defer rows.Close()
	ret := map[int][]*DbCarDriver{}
	for rows.Next() {
		d := DbCarDriver{}
		var eventEntryId int
		err := rows.Scan(&d.ID, &d.CarEntryId, &d.DriverId, &d.Name, &d.Initials,
			&d.AbbrevName, &d.IRating, &d.LicenseLevel, &d.LicenseSubLevel,
			&d.LicenseString, &eventEntryId,
		)
		if err != nil {
			log.Printf("Error scaning driver: %v\n", err)
		}

		if _, ok := ret[eventEntryId]; !ok {
			ret[eventEntryId] = []*DbCarDriver{}
		}
		ret[eventEntryId] = append(ret[eventEntryId], &d)

	}
	return ret, nil
}

//nolint:whitespace // editor/linter issue
func GetDriversByTeam(
	ctx context.Context,
	q internal.Querier,
	teamIds []int,
) (map[int][]*DbCarDriver, error) {
	rows, err := q.Query(ctx, `
	select d.id, d.c_car_entry_id, d.driver_id, d.name, d.initials, 
	  d.abbrev_name, d.irating, d.lic_level,d.lic_sub_level,
	  d.lic_string, e.id
	from c_car_driver d 
	join c_car_entry e on d.c_car_entry_id = e.id
	join c_car_team t on e.id = t.c_car_entry_id
	where t.id = any($1) order by d.name asc`, teamIds)
	if err != nil {
		log.Printf("error reading drivers: %v", err)
		return map[int][]*DbCarDriver{}, err
	}
	defer rows.Close()
	ret := map[int][]*DbCarDriver{}
	for rows.Next() {
		d := DbCarDriver{}
		var teamId int
		err := rows.Scan(&d.ID, &d.CarEntryId, &d.DriverId, &d.Name, &d.Initials,
			&d.AbbrevName, &d.IRating, &d.LicenseLevel, &d.LicenseSubLevel,
			&d.LicenseString, &teamId,
		)
		if err != nil {
			log.Printf("Error scaning driver: %v\n", err)
		}

		if _, ok := ret[teamId]; !ok {
			ret[teamId] = []*DbCarDriver{}
		}
		ret[teamId] = append(ret[teamId], &d)

	}
	return ret, nil
}

// little helper
//
//nolint:unused // by design
const selector = string(`
select d.id, d.c_car_entry_id, d.driver_id, d.name, d.initials, 
	  d.abbrev_name, d.irating, d.lic_level,d.lic_sub_level,d.lic_string, e.id
	from c_car_driver d 
`)

//nolint:unused // by design
func scan(e *DbCarDriver, rows pgx.Rows) error {
	err := rows.Scan(&e.ID, &e.CarEntryId, &e.DriverId, &e.Name, &e.Initials,
		&e.AbbrevName, &e.IRating, &e.LicenseLevel, &e.LicenseSubLevel, &e.LicenseString,
	)

	return err
}
