package events

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DbEvent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string
	Info        struct {
		TrackId           int    `json:"trackId"`
		EventTime         string `json:"eventTime"`
		RaceLoggerVersion string `json:"raceLoggerVersion"`
		TeamRacing        int    `json:"teamRacing"` // 0: false
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

func GetALl(pool *pgxpool.Pool) ([]DbEvent, error) {
	rows, err := pool.Query(context.Background(), selector)
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

func GetById(pool *pgxpool.Pool, id int) (DbEvent, error) {

	res, err := GetByIds(pool, []int{id})
	if err != nil {
		return DbEvent{}, err
	}
	return res[0], nil

}

func GetEventsByTrackIds(pool *pgxpool.Pool, trackIds []int) (map[int][]*DbEvent, error) {

	rows, err := pool.Query(context.Background(), fmt.Sprintf("%s where (data->'info'->'trackId')::integer=any($1)", selector), trackIds)
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
const selector = string("select id,name,event_key,coalesce(description,''),data->'info',data->'manifests', data->'replayInfo' from event ")

func scan(e *DbEvent, rows pgx.Rows) error {
	return rows.Scan(&e.ID, &e.Name, &e.Key, &e.Description, &e.Info, &e.Manifests, &e.ReplayInfo)
}
