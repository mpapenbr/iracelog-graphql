package car

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbCar struct {
	ID            int     `json:"id"`
	EventId       int     `json:"eventId"`
	CarId         int     `json:"carId"`
	Name          string  `json:"name"`
	NameShort     string  `json:"nameShort"`
	FuelPct       float64 `json:"fuelPct"`
	PowerAdjust   float64 `json:"powerAdjust"`
	WeightPenalty float64 `json:"weightPenalty"`
	DryTireSets   int     `json:"dryTireSets"`
}

func GetEventCars(pool *pgxpool.Pool, eventIDs []int) (map[int][]*DbCar, error) {
	rows, err := pool.Query(context.Background(), `
	select id,event_id, car_id, name, name_short, 
	  fuel_pct, power_adjust, weight_penalty, dry_tire_sets
	from c_car
	where event_id = any($1) order by name asc`, eventIDs)
	if err != nil {
		log.Printf("error reading cars: %v", err)
		return map[int][]*DbCar{}, err
	}
	defer rows.Close()
	ret := map[int][]*DbCar{}
	for rows.Next() {
		d := DbCar{}
		err := rows.Scan(&d.ID, &d.EventId, &d.CarId, &d.Name, &d.NameShort,
			&d.FuelPct, &d.PowerAdjust, &d.WeightPenalty, &d.DryTireSets)
		if err != nil {
			log.Printf("Error scanning c_car: %v\n", err)
		}
		if _, ok := ret[d.EventId]; !ok {
			ret[d.EventId] = []*DbCar{}
		}
		ret[d.EventId] = append(ret[d.EventId], &d)
	}
	return ret, nil
}

//nolint:whitespace // editor/linter issue
func GetEventEntryCars(
	pool *pgxpool.Pool,
	eventEntryIDs []int,
) (map[int]*DbCar, error) {
	rows, err := pool.Query(context.Background(), `
	select 
	c.id,
	c.event_id,
	c.car_id,
	c.name,
	c.name_short,
	c.fuel_pct,
	c.power_adjust,
	c.weight_penalty, 
	c.dry_tire_sets,
	ce.id
	from c_car c join c_car_entry ce on ce.c_car_id = c.id
	where ce.id = any($1)`, eventEntryIDs)
	if err != nil {
		log.Printf("error reading cars: %v", err)
		return map[int]*DbCar{}, err
	}
	defer rows.Close()
	ret := map[int]*DbCar{}
	for rows.Next() {
		d := DbCar{}
		var ceId int
		err := rows.Scan(&d.ID, &d.EventId, &d.CarId, &d.Name, &d.NameShort,
			&d.FuelPct, &d.PowerAdjust, &d.WeightPenalty, &d.DryTireSets, &ceId)
		if err != nil {
			log.Printf("Error scanning c_car: %v\n", err)
		}
		ret[ceId] = &d
	}
	return ret, nil
}
