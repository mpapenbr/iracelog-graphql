package tracks

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
)

type DbTrack struct {
	ID   int `json:"id"`
	Data struct {
		Name      string  `json:"trackDisplayName"`
		ShortName string  `json:"trackDisplayShortName"`
		Length    float64 `json:"trackLength"`
	}
}

func GetALl(pool *pgxpool.Pool) ([]DbTrack, error) {
	rows, err := pool.Query(context.Background(), "select id,data from track")
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []DbTrack{}, err
	}
	defer rows.Close()
	var ret []DbTrack
	for rows.Next() {
		t := DbTrack{}
		v, _ := rows.Values()
		log.Printf("%v\n", v)

		err = rows.Scan(&t.ID, &t.Data)
		if err != nil {
			log.Printf("Error scaning Track: %v\n", err)
		}

		log.Printf("%v\n", t)
		ret = append(ret, t)
	}
	return ret, nil
}

func GetByIds(pool *pgxpool.Pool, ids []int) ([]DbTrack, error) {
	arg := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	log.Printf("GetByIds: %v\n", arg)
	rows, err := pool.Query(context.Background(), "select id,data from track where id in ($1)", arg)
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []DbTrack{}, err
	}
	defer rows.Close()
	var ret []DbTrack
	for rows.Next() {
		t := DbTrack{}
		v, _ := rows.Values()
		log.Printf("%v\n", v)

		err = rows.Scan(&t.ID, &t.Data)
		if err != nil {
			log.Printf("Error scaning Track: %v\n", err)
		}

		log.Printf("%v\n", t)
		ret = append(ret, t)
	}
	return ret, nil
}

func GetById(id int) *DbTrack {
	t := DbTrack{}
	err := database.DbPool.QueryRow(context.Background(), "select id,data from track where id=$1", id).Scan(&t.ID, &t.Data)
	if err != nil {
		log.Printf("error reading track: %v", err)
		return nil
	}
	return &t
}
