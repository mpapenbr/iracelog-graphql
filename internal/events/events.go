package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mpapenbr/iracelog-graphql/internal"
)

type DbEvent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string
	RecordStamp time.Time `json:"recordStamp"`
	Info        struct {
		TrackId           int    `json:"trackId"`
		EventTime         string `json:"eventTime"`
		RaceloggerVersion string `json:"raceloggerVersion"`
		TeamRacing        int    `json:"teamRacing"` // 0: false
		MultiClass        bool   `json:"multiClass"`
		NumCarTypes       int    `json:"numCarTypes"`
		NumCarClasses     int    `json:"numCarClasses"`
		IrSessionId       int    `json:"irSessionId"`
		Sessions          []struct {
			Num  int    `json:"num"`
			Name string `json:"name"`
		}
	}
	Manifests struct {
		Car     []string `json:"car"`
		Pit     []string `json:"pit"`
		Message []string `json:"message"`
		Session []string `json:"session"`
	}
	ReplayInfo struct {
		MinTimestamp   float64 `json:"minTimestamp"`
		MinSessionTime float64 `json:"minSessionTime"`
		MaxSessionTime float64 `json:"maxSessionTime"`
	}
}

func GetALl(pool *pgxpool.Pool, pageable internal.DbPageable) ([]DbEvent, error) {
	query := internal.HandlePageableArgs(selector, pageable)

	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("error reading events: %v", err)
		return []DbEvent{}, err
	}
	defer rows.Close()
	var ret []DbEvent
	for rows.Next() {
		e := DbEvent{}
		// log.Printf("%v\n", rows.RawValues())

		err = scan(&e, rows)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}

		// log.Printf("%v\n", e)
		ret = append(ret, e)
	}
	return ret, nil
}

func GetByIds(pool *pgxpool.Pool, ids []int) ([]DbEvent, error) {

	rows, err := pool.Query(context.Background(), fmt.Sprintf("%s where id=any($1)", selector), ids)
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []DbEvent{}, err
	}
	defer rows.Close()
	var ret []DbEvent
	for rows.Next() {
		e := DbEvent{}

		err = scan(&e, rows)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}

		ret = append(ret, e)
	}
	return ret, nil
}

/*
note: currently only pageable.sort is processed.
Discussion: how should limit/offset be interpreted?
We can't put it on the query as this would limit/offset the overall data.
So we have to process it "manually" for each event, which yields the next question:
should offset apply only for those tracks having more than offset events?
consider a track with 2 and another with 10 events and a query with offset 5
*/
func GetEventsByTrackIds(pool *pgxpool.Pool, trackIds []int, pageable internal.DbPageable) (map[int][]*DbEvent, error) {
	query := internal.HandlePageableArgs(fmt.Sprintf("%s where (data->'info'->'trackId')::integer=any($1)", selector), pageable)
	rows, err := pool.Query(context.Background(), query, trackIds)
	if err != nil {
		log.Printf("error reading ids for trackId: %v", err)
		return map[int][]*DbEvent{}, err
	}
	defer rows.Close()
	ret := map[int][]*DbEvent{}
	for rows.Next() {

		var e DbEvent
		err = scan(&e, rows)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}
		val, ok := ret[e.Info.TrackId]
		if !ok {
			val = []*DbEvent{}
		}
		val = append(val, &e)
		ret[e.Info.TrackId] = val
	}

	return ret, nil
}

// little helper
const selector = string("select id,name,event_key,coalesce(description,''),record_stamp, data->'info',data->'manifests', data->'replayInfo' from event ")

func scan(e *DbEvent, rows pgx.Rows) error {
	return rows.Scan(&e.ID, &e.Name, &e.Key, &e.Description, &e.RecordStamp, &e.Info, &e.Manifests, &e.ReplayInfo)
}
