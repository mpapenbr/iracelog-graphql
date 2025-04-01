//nolint:all // code currently not used
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
type DBAnalysis struct {
	ID      int `json:"id"`
	EventID int `json:"eventID"`
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

type DBTeamInEvent struct {
	Name     string `json:"name"`
	CarNum   string `json:"carNum"`
	CarClass string `json:"carClass"`
	Drivers  []struct {
		DriverName string `json:"driverName"`
	}
}

func GetAnalysisForEvent(pool *pgxpool.Pool, eventID int) (*DBAnalysis, error) {
	var data DBAnalysis
	pool.QueryRow(context.Background(), "select id,data from analysis where event_id=$1", eventID).Scan(&data.ID, &data)
	return &data, nil
}

func GetAnalysisForEvents(pool *pgxpool.Pool, eventIDs []int) ([]DBAnalysis, error) {
	rows, err := pool.Query(context.Background(), "select id,event_id,data->'carInfo' from analysis where event_id=any($1)", eventIDs)
	if err != nil {
		log.Printf("error reading analysis: %v", err)
		return []DBAnalysis{}, err
	}
	defer rows.Close()
	ret := []DBAnalysis{}

	for rows.Next() {
		var dba DBAnalysis
		err = rows.Scan(&dba.ID, &dba.EventID, &dba.CarInfo)
		ret = append(ret, dba)
	}
	return ret, err
}

func SearchTeams(pool *pgxpool.Pool, arg string) ([]DBTeamSummary, error) {
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
		return []DBTeamSummary{}, err
	}
	defer rows.Close()
	lookup := map[string]DBTeamSummary{}

	for rows.Next() {
		var dName string
		var carInfo CarInfo
		var eventID int
		err = rows.Scan(&eventID, &dName, &carInfo)
		if err != nil {
			log.Printf("Error scaning result: %v\n", err)
		}
		val, ok := lookup[dName]
		if !ok {
			val = DBTeamSummary{Name: dName}
		}
		val.CarNum = append(val.CarNum, carInfo.CarNum)
		val.CarClass = append(val.CarClass, carInfo.CarClass)
		val.EventIDs = append(val.EventIDs, eventID)
		for _, d := range carInfo.Drivers {
			val.Drivers = append(val.Drivers, d.DriverName)
		}
		lookup[dName] = val

	}

	for k, v := range lookup {
		v.EventIDs = unique(v.EventIDs)
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
	ret := make([]DBTeamSummary, len(keys))
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
