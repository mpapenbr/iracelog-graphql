package analysis

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/slices"
)

type DbDriverSummary struct {
	Name     string   `json:"name"`
	CarNum   []string `json:"carNum"`
	CarClass []string `json:"carClass"`
	Teams    []string `json:"drivers"`
	EventIds []int    `json:"eventId"`
}

func SearchDrivers(pool *pgxpool.Pool, arg string) ([]DbDriverSummary, error) {

	rows, err := pool.Query(context.Background(), fmt.Sprintf(
		`select s.event_id,dInfo->>'driverName' as driverName, tInfo from 
	(select
	a.event_id,
	jsonb_path_query(a.data->'carInfo', '$[*] ?(@.drivers[*].driverName like_regex "%s")') as tInfo,	
	jsonb_path_query(a.data->'carInfo', '$[*].drivers ?(@.driverName like_regex "%s")') as dInfo	
	from analysis a 
	where a.data @? '$.carInfo[*].drivers ?(@.driverName like_regex "%s")'
	) s
	
	`, arg, arg, arg))
	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return []DbDriverSummary{}, err
	}
	defer rows.Close()
	lookup := map[string]DbDriverSummary{}

	for rows.Next() {
		var dName string
		var carInfo CarInfo
		var eventId int
		err = rows.Scan(&eventId, &dName, &carInfo)
		if err != nil {
			log.Printf("Error scaning result: %v\n", err)
		}
		val, ok := lookup[dName]
		if !ok {
			val = DbDriverSummary{Name: dName}
		}
		val.CarNum = append(val.CarNum, carInfo.CarNum)
		val.CarClass = append(val.CarClass, carInfo.CarClass)
		val.EventIds = append(val.EventIds, eventId)

		val.Teams = append(val.Teams, carInfo.Name)

		lookup[dName] = val

	}

	for k, v := range lookup {
		v.EventIds = unique(v.EventIds)
		v.Teams = unique(v.Teams)
		v.CarNum = unique(v.CarNum)
		v.CarClass = unique(v.CarClass)
		sort.Strings(v.Teams)
		sort.Strings(v.CarNum)
		sort.Strings(v.CarClass)
		lookup[k] = v
	}
	keys := make([]string, 0, len(lookup))
	for k := range lookup {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := make([]DbDriverSummary, len(keys))
	for i, v := range keys {
		ret[i] = lookup[v]
	}
	return ret, nil
}

/*
teams contains the exact team names for which the drivers should be collected.
*/
func SearchDriversInTeams(pool *pgxpool.Pool, teams []string) (map[string][]DbDriverSummary, error) {
	work := make([]string, len(teams))
	for i, t := range teams {
		work[i] = fmt.Sprintf("('%s')", t)
	}
	strings.Join(work, ",")
	rows, err := pool.Query(context.Background(), fmt.Sprintf(
		`select s.event_id,s.teamName, tInfo from 
	(select
	a.event_id,
	jsonb_path_query(a.data, '$.carInfo[*] ?(@.name == $myArg)', jsonb_build_object('myArg', args.arg) ) as tInfo,
	args.arg as teamName
	from analysis a cross join (select * from (values %s) as b(arg)) args
	where jsonb_path_exists(a.data, '$.carInfo[*] ?(@.name == $myArg)', jsonb_build_object('myArg', args.arg) )
	) s
	order by teamName
	`, strings.Join(work, ",")))
	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return map[string][]DbDriverSummary{}, err
	}
	defer rows.Close()
	teamLookup := map[string][]DbDriverSummary{}
	driverLookup := map[string]DbDriverSummary{}

	for rows.Next() {
		var teamName string
		var carInfo CarInfo
		var eventId int
		err = rows.Scan(&eventId, &teamName, &carInfo)
		if err != nil {
			log.Printf("Error scaning result: %v\n", err)
		}
		teamEntry, ok := teamLookup[teamName]
		if !ok {
			teamEntry = []DbDriverSummary{}
		}
		for _, d := range carInfo.Drivers {
			driverEntry, ok := driverLookup[d.DriverName]
			if !ok {
				driverEntry = DbDriverSummary{Name: d.DriverName}
			}
			driverEntry.CarNum = append(driverEntry.CarNum, carInfo.CarNum)
			driverEntry.CarClass = append(driverEntry.CarClass, carInfo.CarClass)
			driverEntry.EventIds = append(driverEntry.EventIds, eventId)
			driverEntry.Teams = append(driverEntry.Teams, teamName)
			driverLookup[d.DriverName] = driverEntry

		}
		teamLookup[teamName] = teamEntry

	}

	for k, v := range driverLookup {
		v.EventIds = unique(v.EventIds)
		v.Teams = unique(v.Teams)
		v.CarNum = unique(v.CarNum)
		v.CarClass = unique(v.CarClass)
		sort.Strings(v.Teams)
		sort.Strings(v.CarNum)
		sort.Strings(v.CarClass)
		driverLookup[k] = v
	}

	// TODO: refactor to generic util func
	keys := make([]string, 0, len(teamLookup))
	for k := range teamLookup {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for k := range teamLookup {
		for _, d := range driverLookup {
			if slices.Contains(d.Teams, k) {
				val := teamLookup[k]
				val = append(val, d)
				teamLookup[k] = val
			}
		}
	}
	return teamLookup, nil
}

func CollectEventIdsForDriver(pool *pgxpool.Pool, driverName string) ([]int, error) {
	rows, err := pool.Query(context.Background(), fmt.Sprintf(
		`select a.event_id from analysis a where a.data @? '$.carInfo[*].drivers ?(@.driverName =="%s")'`, driverName))

	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return []int{}, err
	}
	defer rows.Close()

	ret := []int{}
	for rows.Next() {
		var eventId int
		err = rows.Scan(&eventId)
		if err != nil {
			log.Printf("Error scaning result: %v\n", err)
		}
		ret = append(ret, eventId)
	}
	return ret, nil
}

func CollectEventIdsForDrivers(pool *pgxpool.Pool, driverNames []string) (map[string][]int, error) {
	work := make([]string, len(driverNames))
	for i, t := range driverNames {
		work[i] = fmt.Sprintf("('%s')", t)
	}
	strings.Join(work, ",")
	rows, err := pool.Query(context.Background(), fmt.Sprintf(
		`select a.event_id, args.arg from analysis a cross join (select * from (values %s) as b(arg)) args
		where jsonb_path_exists(a.data, '$.carInfo[*] ?(@.drivers.driverName == $myArg)', jsonb_build_object('myArg', args.arg))
		`, strings.Join(work, ",")))

	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return map[string][]int{}, err
	}
	defer rows.Close()

	ret := map[string][]int{}
	for rows.Next() {
		var eventId int
		var driver string
		err = rows.Scan(&eventId, &driver)
		v, ok := ret[driver]
		if !ok {
			v = []int{}
		}
		v = append(v, eventId)
		ret[driver] = v
	}
	return ret, nil
}
