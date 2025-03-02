package events

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mpapenbr/iracelog-graphql/internal"
)

type DbEvent struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Key                  string `json:"key"`
	Description          string
	EventTime            time.Time `json:"eventTime"`
	RaceloggerVersion    string    `json:"raceloggerVersion"`
	TeamRacing           bool      `json:"teamRacing"`
	MultiClass           bool      `json:"multiClass"`
	NumCarTypes          int       `json:"numCarTypes"`
	NumCarClasses        int       `json:"numCarClasses"`
	IrSessionId          int       `json:"irSessionId"`
	TrackId              int       `json:"trackId"`
	PitSpeed             float64   `json:"pitSpeed"`
	ReplayMinTimestamp   time.Time `json:"replayMinTimestamp"`
	ReplayMinSessionTime float64   `json:"replayMinSessionTime"`
	ReplayMaxSessionTime float64   `json:"replayMaxSessionTime"`
	Sessions             []Session `json:"sessions"`
}

//nolint:tagliatelle // json is that way
type Session struct {
	Num         int    `json:"num"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	SessionTime int    `json:"session_time"`
	Laps        int    `json:"laps"`
}

type EventSearchKeys struct {
	Name   string
	Car    string
	Track  string
	Driver string
	Team   string
}

func GetALl(pool *pgxpool.Pool, pageable internal.DbPageable) ([]*DbEvent, error) {
	query := internal.HandlePageableArgs(selector, pageable)

	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("error reading events: %v", err)
		return []*DbEvent{}, err
	}
	defer rows.Close()
	var ret []*DbEvent
	for rows.Next() {
		e := DbEvent{}
		err = scan(&e, rows)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}
		ret = append(ret, &e)
	}
	return ret, nil
}

func GetByIds(pool *pgxpool.Pool, ids []int) ([]*DbEvent, error) {
	rows, err := pool.Query(context.Background(),
		fmt.Sprintf("%s where id=any($1)", selector), ids)
	if err != nil {
		log.Printf("error reading tracks: %v", err)
		return []*DbEvent{}, err
	}
	defer rows.Close()
	var ret []*DbEvent
	for rows.Next() {
		e := DbEvent{}
		err = scan(&e, rows)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}
		ret = append(ret, &e)
	}
	return ret, nil
}

/*
note: currently only pageable.sort is processed.
Discussion: how should limit/offset be interpreted?
We can't put it on the query as this would limit/offset the overall data.
So we have to process it "manually" for each event, which yields the next question:
should offset apply only for those tracks having more than offset events?
consider a track with 2 and another with 10 events and a query with offset 5
*/
//nolint:whitespace // editor/linter issue
func GetEventsByTrackIds(
	pool *pgxpool.Pool,
	trackIds []int,
	pageable internal.DbPageable,
) (map[int][]*DbEvent, error) {
	query := internal.HandlePageableArgs(
		fmt.Sprintf("%s where track_id=any($1)", selector), pageable)
	rows, err := pool.Query(context.Background(), query, trackIds)
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
		val, ok := ret[e.TrackId]
		if !ok {
			val = []*DbEvent{}
		}
		val = append(val, &e)
		ret[e.TrackId] = val
	}

	return ret, nil
}

//nolint:lll,whitespace // sql redability
func SimpleEventSearch(
	pool *pgxpool.Pool,
	searchArg string,
	pageable internal.DbPageable,
) ([]*DbEvent, error) {
	query := internal.HandlePageableArgs(fmt.Sprintf(
		`%s
WHERE name ilike $1
OR    description ilike $1
OR    track_id in (select id from track where name ilike $1)
OR id in (select event_id from c_car where name ilike $1)
OR id in (select e.event_id from c_car_entry e join c_car_team t on t.c_car_entry_id=e.id and t.name ilike $1)
OR id in (select e.event_id from c_car_entry e join c_car_driver d on d.c_car_entry_id=e.id and d.name ilike $1)
		`, selector), pageable)

	rows, err := pool.Query(context.Background(), query, fmt.Sprintf("%%%s%%", searchArg))
	if err != nil {
		log.Printf("error reading ids for searchArg: %v", err)
		return []*DbEvent{}, err
	}
	defer rows.Close()
	var ret []*DbEvent
	for rows.Next() {
		e := DbEvent{}
		err = scan(&e, rows)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}
		ret = append(ret, &e)
	}
	return ret, nil
}

//nolint:lll,funlen,whitespace // sql redability
func AdvancedEventSearch(
	pool *pgxpool.Pool,
	search *EventSearchKeys,
	pageable internal.DbPageable,
) ([]*DbEvent, error) {
	type paramType struct {
		Selector string
		Param    *EventSearchKeys
	}
	param := paramType{Selector: selector, Param: search}
	//nolint:lll // sql redability
	tmpl, err := template.New("sql").Parse(`
	{{ .Selector }}
	WHERE
	{{if .Param.Name }} name ilike '%{{ .Param.Name }}%' {{ else }} true {{ end }}
	
	{{- if .Param.Track }} 
	AND track_id in (select id from track where name ilike '%{{ .Param.Track }}%') 
	{{ end }}
	 
	{{- if .Param.Car }}
	AND id in (select event_id from c_car where name ilike '%{{ .Param.Car }}%') 
	{{ end }}

	{{- if .Param.Team }}
	AND id in (select e.event_id from c_car_entry e join c_car_team t on t.c_car_entry_id=e.id and t.name ilike '%{{ .Param.Team }}%')
	{{ end }}
	
	{{- if .Param.Driver }}
	AND id in (select e.event_id from c_car_entry e join c_car_driver d on d.c_car_entry_id=e.id and d.name ilike '%{{ .Param.Driver }}%')
	{{ end }}
	
	`)
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, param)
	if err != nil {
		return nil, err
	}
	qString := tpl.String()
	query := internal.HandlePageableArgs(qString, pageable)

	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("error reading ids for searchArg: %v", err)
		return []*DbEvent{}, err
	}
	defer rows.Close()
	var ret []*DbEvent
	for rows.Next() {
		e := DbEvent{}
		err = scan(&e, rows)
		if err != nil {
			log.Printf("Error scaning Event: %v\n", err)
		}
		ret = append(ret, &e)
	}
	return ret, nil
}

// little helper
const selector = string(`
select id,name,event_key,description,event_time,
racelogger_version,team_racing, multi_class, num_car_types,
num_car_classes,ir_session_id, track_id, pit_speed,
replay_min_timestamp, replay_min_session_time,replay_max_session_time, sessions
from event 
`)

func scan(e *DbEvent, rows pgx.Rows) error {
	err := rows.Scan(&e.ID, &e.Name, &e.Key, &e.Description, &e.EventTime,
		&e.RaceloggerVersion, &e.TeamRacing, &e.MultiClass, &e.NumCarTypes,
		&e.NumCarClasses, &e.IrSessionId, &e.TrackId, &e.PitSpeed,
		&e.ReplayMinTimestamp, &e.ReplayMinSessionTime, &e.ReplayMaxSessionTime,
		&e.Sessions,
	)

	return err
}
