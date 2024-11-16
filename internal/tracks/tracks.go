package tracks

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mpapenbr/iracelog-graphql/internal"
)

type DbTrack struct {
	ID int `json:"id"`

	Name          string   `json:"name"`
	ShortName     string   `json:"shortName"`
	Config        string   `json:"config"`
	Length        float64  `json:"length"`
	PitSpeed      float64  `json:"pitSpeed"`
	PitExit       float64  `json:"pitExit"`
	PitEntry      float64  `json:"pitEntry"`
	PitLaneLength float64  `json:"pitLaneLength"`
	Sectors       []Sector `json:"sectors"`
}
type Sector struct {
	Num      int     `json:"num"`
	StartPct float64 `json:"start_pct"`
}

// github.com/cweill/gotests/gotests@v1.6.0
func GetALl(pool *pgxpool.Pool, pageable internal.DbPageable) ([]*DbTrack, error) {
	query := internal.HandlePageableArgs(selector, pageable)
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []*DbTrack{}, err
	}
	defer rows.Close()
	ret := []*DbTrack{}
	for rows.Next() {
		// t := DbTrack{}
		// v, _ := rows.Values()
		// log.Printf("%v\n", v)

		t, err := scanRow(rows)
		if err != nil {
			log.Printf("Error scaning Track: %v\n", err)
		}

		// log.Printf("%v\n", t)
		ret = append(ret, t)
	}
	return ret, nil
}

func GetByIds(pool *pgxpool.Pool, ids []int) ([]*DbTrack, error) {
	rows, err := pool.Query(context.Background(), fmt.Sprintf("%s where id=any($1)", selector), ids)
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []*DbTrack{}, err
	}
	defer rows.Close()
	var ret []*DbTrack
	for rows.Next() {
		t, err := scanRow(rows)
		if err != nil {
			log.Printf("Error scaning Track: %v\n", err)
		}

		ret = append(ret, t)
	}
	return ret, nil
}

// little helper
const selector = string("select id,name,short_name,config,track_length,sectors,pit_speed,pit_entry,pit_exit,pit_lane_length from track")

func scanRow(row pgx.Row) (*DbTrack, error) {
	var sectors []Sector
	var item DbTrack

	if err := row.Scan(&item.ID, &item.Name, &item.ShortName, &item.Config,
		&item.Length, &sectors, &item.PitSpeed,
		&item.PitEntry, &item.PitExit, &item.PitLaneLength); err != nil {
		return nil, err
	}
	item.Sectors = make([]Sector, len(sectors))
	for i := range sectors {
		item.Sectors[i] = sectors[i]
	}

	return &item, nil
}
