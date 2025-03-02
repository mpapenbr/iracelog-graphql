package team

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbCarTeam struct {
	ID     int    `json:"id"`
	TeamId int    `json:"teamId"`
	Name   string `json:"name"`
}

//nolint:whitespace // editor/linter issue
func GetTeamsByEventEntry(
	pool *pgxpool.Pool,
	eventEntryIDs []int,
) (map[int]*DbCarTeam, error) {
	rows, err := pool.Query(context.Background(), `
	select 
	t.id,	
	t.name,	
	t.team_id,	
	ce.id
	from c_car_team t join c_car_entry ce on ce.id = t.c_car_entry_id
	where ce.id = any($1)`, eventEntryIDs)
	if err != nil {
		log.Printf("error reading car_team: %v", err)
		return map[int]*DbCarTeam{}, err
	}
	defer rows.Close()
	ret := map[int]*DbCarTeam{}
	for rows.Next() {
		d := DbCarTeam{}
		var ceId int
		err := rows.Scan(&d.ID, &d.Name, &d.TeamId, &ceId)
		if err != nil {
			log.Printf("Error scanning c_car_team: %v\n", err)
		}
		ret[ceId] = &d
	}
	return ret, nil
}
