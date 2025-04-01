//nolint:all // code currently not used
package analysis

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slices"
)

type DBTeamSummary struct {
	Name     string   `json:"name"`
	CarNum   []string `json:"carNum"`
	CarClass []string `json:"carClass"`
	Drivers  []string `json:"drivers"`
	EventIDs []int    `json:"eventID"`
}

/*
teams contains the exact team names for which the drivers should be collected.
*/
func SearchTeamsForDrivers(pool *pgxpool.Pool, drivers []string) (map[string][]DBTeamSummary, error) {
	work := make([]string, len(drivers))
	for i, t := range drivers {
		work[i] = fmt.Sprintf("('%s')", t)
	}
	strings.Join(work, ",")
	rows, err := pool.Query(context.Background(), fmt.Sprintf(
		`select s.event_id,s.driverName, tInfo from 
	(select
	a.event_id,
	jsonb_path_query(a.data, '$.carInfo[*] ?(@.drivers.driverName == $myArg)', jsonb_build_object('myArg', args.arg) ) as tInfo,
	args.arg as driverName
	from analysis a cross join (select * from (values %s) as b(arg)) args
	where jsonb_path_exists(a.data, '$.carInfo[*].drivers[*] ?(@.driverName == $myArg)', jsonb_build_object('myArg', args.arg) )
	) s
	
	`, strings.Join(work, ",")))
	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return map[string][]DBTeamSummary{}, err
	}
	defer rows.Close()
	teamLookup := map[string]DBTeamSummary{}
	driverLookup := map[string][]DBTeamSummary{}

	for rows.Next() {
		var driverName string
		var carInfo CarInfo
		var eventID int
		err = rows.Scan(&eventID, &driverName, &carInfo)
		if err != nil {
			log.Printf("Error scaning result: %v\n", err)
		}
		driverEntry, ok := driverLookup[driverName]
		if !ok {
			driverEntry = []DBTeamSummary{}
			driverLookup[driverName] = driverEntry
		}
		teamEntry, ok := teamLookup[carInfo.Name]
		if !ok {
			teamEntry = DBTeamSummary{Name: carInfo.Name}
			// driverEntry = append(driverEntry, teamEntry)
			// driverLookup[driverName] = driverEntry
		}

		teamEntry.CarNum = append(teamEntry.CarNum, carInfo.CarNum)
		teamEntry.CarClass = append(teamEntry.CarClass, carInfo.CarClass)
		teamEntry.EventIDs = append(teamEntry.EventIDs, eventID)
		teamEntry.Drivers = append(teamEntry.Drivers, driverName)

		teamLookup[carInfo.Name] = teamEntry

	}

	for k, v := range teamLookup {
		v.EventIDs = unique(v.EventIDs)

		v.CarNum = unique(v.CarNum)
		v.CarClass = unique(v.CarClass)

		sort.Strings(v.CarNum)
		sort.Strings(v.CarClass)
		teamLookup[k] = v
	}

	// TODO: refactor to generic util func
	keys := make([]string, 0, len(driverLookup))
	for k := range driverLookup {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for k := range driverLookup {
		for _, t := range teamLookup {
			if slices.Contains(t.Drivers, k) {
				val := driverLookup[k]
				val = append(val, t)
				driverLookup[k] = val
			}
		}
	}
	return driverLookup, nil
}

func CollectEventIDsForTeams(pool *pgxpool.Pool, teamNames []string) (map[string][]int, error) {
	work := make([]string, len(teamNames))
	for i, t := range teamNames {
		work[i] = fmt.Sprintf("('%s')", t)
	}
	strings.Join(work, ",")
	rows, err := pool.Query(context.Background(), fmt.Sprintf(
		`select a.event_id, args.arg from analysis a cross join (select * from (values %s) as b(arg)) args
		where jsonb_path_exists(a.data, '$.carInfo[*] ?(@.name == $myArg)', jsonb_build_object('myArg', args.arg))
		`, strings.Join(work, ",")))
	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return map[string][]int{}, err
	}
	defer rows.Close()

	ret := map[string][]int{}
	for rows.Next() {
		var eventID int
		var team string
		err = rows.Scan(&eventID, &team)
		v, ok := ret[team]
		if !ok {
			v = []int{}
		}
		v = append(v, eventID)
		ret[team] = v
	}
	return ret, nil
}
