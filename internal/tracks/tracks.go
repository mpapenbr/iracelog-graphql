package tracks

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mpapenbr/iracelog-graphql/internal"
	database "github.com/mpapenbr/iracelog-graphql/internal/pkg/db/postgres"
)

type DbTrack struct {
	ID   int `json:"id"`
	Data struct {
		Name      string  `json:"trackDisplayName"`
		ShortName string  `json:"trackDisplayShortName"`
		Config    string  `json:"trackConfigName"`
		Length    float64 `json:"trackLength"`
		Pit       struct {
			Exit       float64 `json:"exit"`
			Entry      float64 `json:"entry"`
			LaneLength float64 `json:"laneLength"`
		} `json:"pit"`
		Sectors []struct {
			SectorNum      int     `json:"sectorNum"`
			SectorStartpct float64 `json:"sectorStartPct"`
		} `json:"sectors"`
	}
}

func GetALl(pool *pgxpool.Pool, pageable internal.DbPageable) ([]DbTrack, error) {
	query := internal.HandlePageableArgs(selector, pageable)
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []DbTrack{}, err
	}
	defer rows.Close()
	var ret []DbTrack
	for rows.Next() {
		t := DbTrack{}
		// v, _ := rows.Values()
		// log.Printf("%v\n", v)

		err = scan(&t, rows)
		if err != nil {
			log.Printf("Error scaning Track: %v\n", err)
		}

		// log.Printf("%v\n", t)
		ret = append(ret, t)
	}
	return ret, nil
}

func GetByIds(pool *pgxpool.Pool, ids []int) ([]DbTrack, error) {

	rows, err := pool.Query(context.Background(), fmt.Sprintf("%s where id=any($1)", selector), ids)
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []DbTrack{}, err
	}
	defer rows.Close()
	var ret []DbTrack
	for rows.Next() {
		t := DbTrack{}

		err = scan(&t, rows)
		if err != nil {
			log.Printf("Error scaning Track: %v\n", err)
		}

		ret = append(ret, t)
	}
	return ret, nil
}

func GetById(id int) *DbTrack {
	t := DbTrack{}
	err := scanRow(&t, database.DbPool.QueryRow(context.Background(), fmt.Sprintf("%s where id=$1", selector), id))
	if err != nil {
		log.Printf("error reading track: %v", err)
		return nil
	}
	return &t
}

// little helper
const selector = string("select id,data from track")

func scan(t *DbTrack, rows pgx.Rows) error {
	return rows.Scan(&t.ID, &t.Data)
}
func scanRow(t *DbTrack, row pgx.Row) error {
	return row.Scan(&t.ID, &t.Data)
}
