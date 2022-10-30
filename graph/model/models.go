package model

import (
	"time"

	"github.com/mpapenbr/iracelog-graphql/internal/analysis"
	"github.com/mpapenbr/iracelog-graphql/internal/events"
)

// contains our own models which are not generated by gqlgen.
// These models are used to glue both worlds (see definitions in internal and generated graphql stuff)
// we may put data into those structs during resolving

type Event struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Key               string    `json:"key"`
	TrackId           int       `json:"trackId"`
	RecordDate        time.Time `json:"recordDate"`
	EventDate         time.Time `json:"eventDate"`
	RaceloggerVersion string    `json:"raceloggerVersion"`
	TeamRacing        bool      `json:"teamRacing"`
	MultiClass        bool      `json:"multiClass"`
	IracingSessionId  int       `json:"iracingSessionId"`
	NumCarClasses     int       `json:"numCarClasses"`
	NumCarTypes       int       `json:"numCarTypes"`

	Track          *Track               `json:"track"`
	DbEvent        *events.DbEvent      `json:"dbEvent"`
	DbAnalysisData *analysis.DbAnalysis `json:"analysisData"`
}

type Track struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ShortName string  `json:"shortName"`
	Length    float64 `json:"length"`
}

type Pageable struct {
	Limit  *int
	Offset *int
}

type EventPageable struct {
	Pageable
	Sort []EventSortArg
}
