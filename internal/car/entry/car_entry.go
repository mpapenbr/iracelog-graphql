package entry

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbCarEntry struct {
	ID        int    `json:"id"`
	EventId   int    `json:"eventId"`
	CarId     int    `json:"carId"`
	CarIdx    int    `json:"carIdx"`
	CarNum    string `json:"carNum"`
	CarNumRaw int    `json:"carNumRaw"`
}

func GetEventEntriesByEventId(pool *pgxpool.Pool, eventIDs []int) (map[int][]*DbCarEntry, error) {
	rows, err := pool.Query(context.Background(), `
	select id,event_id, c_car_id, car_idx, car_number, car_number_raw
	from c_car_entry
	where event_id = any($1) order by car_number_raw asc`, eventIDs)
	if err != nil {
		log.Printf("error reading entries: %v", err)
		return map[int][]*DbCarEntry{}, err
	}
	defer rows.Close()
	ret := map[int][]*DbCarEntry{}
	for rows.Next() {
		d := DbCarEntry{}
		err := rows.Scan(&d.ID, &d.EventId, &d.CarId, &d.CarIdx, &d.CarNum, &d.CarNumRaw)
		if err != nil {
			log.Printf("Error scanning c_car_entry: %v\n", err)
		}

		if _, ok := ret[d.EventId]; !ok {
			ret[d.EventId] = []*DbCarEntry{}
		}
		ret[d.EventId] = append(ret[d.EventId], &d)

	}
	return ret, nil
}

func GetEventEntriesByIds(pool *pgxpool.Pool, ids []int) (map[int]*DbCarEntry, error) {
	rows, err := pool.Query(context.Background(), `
	select id,event_id, c_car_id, car_idx, car_number, car_number_raw
	from c_car_entry
	where id = any($1) order by car_number_raw asc`, ids)
	if err != nil {
		log.Printf("error reading entries: %v", err)
		return map[int]*DbCarEntry{}, err
	}
	defer rows.Close()
	ret := map[int]*DbCarEntry{}
	for rows.Next() {
		d := DbCarEntry{}
		err := rows.Scan(&d.ID, &d.EventId, &d.CarId, &d.CarIdx, &d.CarNum, &d.CarNumRaw)
		if err != nil {
			log.Printf("Error scanning c_car_entry: %v\n", err)
		}

		ret[d.ID] = &d

	}
	return ret, nil
}
