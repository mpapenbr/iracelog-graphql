package tracks

import (
	"context"
	"log"

	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
)

type Track struct {
	ID   int64 `json:"id"`
	Data struct {
		Name      string  `json:"trackDisplayName"`
		ShortName string  `json:"trackDisplayShortName"`
		Length    float64 `json:"trackLength"`
	}
}

func GetALl() []Track {
	rows, err := database.DbPool.Query(context.Background(), "select id,data from track")
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []Track{}
	}
	defer rows.Close()
	var ret []Track
	for rows.Next() {
		t := Track{}
		v, _ := rows.Values()
		log.Printf("%v\n", v)

		err = rows.Scan(&t.ID, &t.Data)
		if err != nil {
			log.Printf("Error scaning Track: %v\n", err)
		}

		log.Printf("%v\n", t)
		ret = append(ret, t)
	}
	return ret
}

func GetById(id int64) *Track {
	t := Track{}
	err := database.DbPool.QueryRow(context.Background(), "select id,data from track where id=$1", id).Scan(&t.ID, &t.Data)
	if err != nil {
		log.Printf("error reading track: %v", err)
		return nil
	}
	return &t
}
