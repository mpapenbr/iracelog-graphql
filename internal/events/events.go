package events

import (
	"context"
	"log"

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
		Sessions          []struct {
			Num  int    `json:"num"`
			Name string `json:"name"`
		}
	}
}

func GetALl(pool *pgxpool.Pool) ([]DbEvent, error) {
	rows, err := pool.Query(context.Background(), "select id,name,coalesce(description,''),data->'info' from event")
	if err != nil {
		log.Printf("error reading events: %v", err)
		return []DbEvent{}, err
	}
	defer rows.Close()
	var ret []DbEvent
	for rows.Next() {
		e := DbEvent{}
		// log.Printf("%v\n", rows.RawValues())

		err = rows.Scan(&e.ID, &e.Name, &e.Description, &e.Info)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}

		log.Printf("%v\n", e)
		ret = append(ret, e)
	}
	return ret, nil
}
