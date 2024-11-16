package analysis

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CarInfo struct {
	Name     string `json:"name"`
	CarNum   string `json:"carNum"`
	CarClass string `json:"carClass"`
	Drivers  []struct {
		SeatTime []struct {
			EnterCarTime float64 // unit: sessionTime
			LeaveCarTime float64 // unit: sessionTime
		}
		DriverName string `json:"driverName"`
	}
}
type DbAnalysis struct {
	ID      int `json:"id"`
	EventId int `json:"eventId"`
	// TODO: Cars
	CarInfo []CarInfo
	CarLaps []struct {
		CarNum string
		Laps   []struct {
			LapNo   int
			LapTime float64
		}
	}
	RaceOrder []string // last computed race order (contain carNums)
}

type DbTeamInEvent struct {
	Name     string `json:"name"`
	CarNum   string `json:"carNum"`
	CarClass string `json:"carClass"`
	Drivers  []struct {
		DriverName string `json:"driverName"`
	}
}

func GetAnalysisForEvent(pool *pgxpool.Pool, eventId int) (*DbAnalysis, error) {
	var data DbAnalysis
	pool.QueryRow(context.Background(), "select id,data from analysis where event_id=$1", eventId).Scan(&data.ID, &data)
	return &data, nil
}

func GetAnalysisForEvents(pool *pgxpool.Pool, eventIds []int) ([]DbAnalysis, error) {
	rows, err := pool.Query(context.Background(), "select id,event_id,data->'carInfo' from analysis where event_id=any($1)", eventIds)
	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return []DbAnalysis{}, err
	}
	defer rows.Close()
	ret := []DbAnalysis{}

	for rows.Next() {
		var dba DbAnalysis
		err = rows.Scan(&dba.ID, &dba.EventId, &dba.CarInfo)
		ret = append(ret, dba)
	}
	return ret, err
}

func SearchTeams(pool *pgxpool.Pool, arg string) ([]DbTeamSummary, error) {
	rows, err := pool.Query(context.Background(), fmt.Sprintf(
		`select s.event_id,tInfo->>'name' as teamName, tInfo from 
	(select
	a.event_id,
	jsonb_path_query(a.data->'carInfo', '$[*] ?(@.name like_regex "%s")') as tInfo	
	from analysis a 
	where a.data @? '$.carInfo[*] ? (@.name like_regex "%s")'
	) s
	order by teamName
	`, arg, arg))
	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return []DbTeamSummary{}, err
	}
	defer rows.Close()
	lookup := map[string]DbTeamSummary{}

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
			val = DbTeamSummary{Name: dName}
		}
		val.CarNum = append(val.CarNum, carInfo.CarNum)
		val.CarClass = append(val.CarClass, carInfo.CarClass)
		val.EventIds = append(val.EventIds, eventId)
		for _, d := range carInfo.Drivers {
			val.Drivers = append(val.Drivers, d.DriverName)
		}
		lookup[dName] = val

	}

	for k, v := range lookup {
		v.EventIds = unique(v.EventIds)
		v.Drivers = unique(v.Drivers)
		v.CarNum = unique(v.CarNum)
		v.CarClass = unique(v.CarClass)
		sort.Strings(v.Drivers)
		sort.Strings(v.CarNum)
		sort.Strings(v.CarClass)
		lookup[k] = v
	}
	keys := make([]string, 0, len(lookup))
	for k := range lookup {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := make([]DbTeamSummary, len(keys))
	for i, v := range keys {
		ret[i] = lookup[v]
	}
	return ret, nil
}

func unique[T comparable](s []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
