package events

import (
	"context"
	"log"

	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
)

type Event struct {
	ID          int64  `json:"id"`
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

func GetALl() []Event {
	rows, err := database.DbPool.Query(context.Background(), "select id,name,coalesce(description,''),data->'info' from event")
	if err != nil {
		log.Printf("error reading events: %v", err)
		return []Event{}
	}
	defer rows.Close()
	var ret []Event
	for rows.Next() {
		e := Event{}
		// log.Printf("%v\n", rows.RawValues())

		err = rows.Scan(&e.ID, &e.Name, &e.Description, &e.Info)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}

		log.Printf("%v\n", e)
		ret = append(ret, e)
	}
	return ret
}
