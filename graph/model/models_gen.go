// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Car struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	NameShort string `json:"nameShort"`
	CarID     int    `json:"carId"`
	// fuel capacity in percent
	FuelPct float64 `json:"fuelPct"`
	// engine power adjustment percent
	PowerAdjust float64 `json:"powerAdjust"`
	// weigth penalty in kg
	WeightPenalty float64 `json:"weightPenalty"`
	// number of dry tire sets
	DryTireSets int `json:"dryTireSets"`
}

// This models a more 'generic' driver with participation in events and teams.
type Driver struct {
	// The driver name used
	Name string `json:"name"`
	// The teams in which the driver was a member
	Teams []*Team `json:"teams"`
	// The events in which the driver participated
	Events []*Event `json:"events"`
	// The car numbers used by this driver
	CarNum []string `json:"carNum"`
	// The car classes used by this driver
	CarClass []string `json:"carClass"`
}

// This models a driver in a concrete event
type EventDriver struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	DriverID        int     `json:"driverId"`
	Initials        *string `json:"initials,omitempty"`
	AbbrevName      *string `json:"abbrevName,omitempty"`
	IRating         *int    `json:"iRating,omitempty"`
	LicenseLevel    *int    `json:"licenseLevel,omitempty"`
	LicenseSubLevel *int    `json:"licenseSubLevel,omitempty"`
	LicenseString   *string `json:"licenseString,omitempty"`
}

// describes an entry in a specific event.
type EventEntry struct {
	ID int `json:"id"`
	// The car data with optional specific restrictions
	Car *Car `json:"car"`
	// The car number for this car
	CarNum *string `json:"carNum,omitempty"`
	// The car number in iRacing raw format
	CarNumRaw *int `json:"carNumRaw,omitempty"`
	// the team running this car
	Team *EventTeam `json:"team,omitempty"`
	// the drivers of this car
	Drivers []*EventDriver `json:"drivers"`
}

type EventSortArg struct {
	Field EventSortField `json:"field"`
	Order *SortOrder     `json:"order,omitempty"`
}

type EventTeam struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	TeamID  int            `json:"teamId"`
	Drivers []*EventDriver `json:"drivers"`
}

type Query struct {
}

// This models a more 'generic' driver with participation in events and teams.
type Team struct {
	Name     string       `json:"name"`
	Drivers  []*Driver    `json:"drivers"`
	CarNum   []string     `json:"carNum"`
	CarClass []string     `json:"carClass"`
	Teams    []*EventTeam `json:"teams"`
	Events   []*Event     `json:"events"`
}

type TrackSortArg struct {
	Field TrackSortField `json:"field"`
	Order *SortOrder     `json:"order,omitempty"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EventSortField string

const (
	EventSortFieldName       EventSortField = "NAME"
	EventSortFieldRecordDate EventSortField = "RECORD_DATE"
	EventSortFieldTrack      EventSortField = "TRACK"
)

var AllEventSortField = []EventSortField{
	EventSortFieldName,
	EventSortFieldRecordDate,
	EventSortFieldTrack,
}

func (e EventSortField) IsValid() bool {
	switch e {
	case EventSortFieldName, EventSortFieldRecordDate, EventSortFieldTrack:
		return true
	}
	return false
}

func (e EventSortField) String() string {
	return string(e)
}

func (e *EventSortField) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventSortField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventSortField", str)
	}
	return nil
}

func (e EventSortField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

var AllSortOrder = []SortOrder{
	SortOrderAsc,
	SortOrderDesc,
}

func (e SortOrder) IsValid() bool {
	switch e {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

func (e SortOrder) String() string {
	return string(e)
}

func (e *SortOrder) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortOrder", str)
	}
	return nil
}

func (e SortOrder) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TrackSortField string

const (
	TrackSortFieldID            TrackSortField = "ID"
	TrackSortFieldName          TrackSortField = "NAME"
	TrackSortFieldShortName     TrackSortField = "SHORT_NAME"
	TrackSortFieldLength        TrackSortField = "LENGTH"
	TrackSortFieldPitlaneLength TrackSortField = "PITLANE_LENGTH"
	TrackSortFieldNumSectors    TrackSortField = "NUM_SECTORS"
)

var AllTrackSortField = []TrackSortField{
	TrackSortFieldID,
	TrackSortFieldName,
	TrackSortFieldShortName,
	TrackSortFieldLength,
	TrackSortFieldPitlaneLength,
	TrackSortFieldNumSectors,
}

func (e TrackSortField) IsValid() bool {
	switch e {
	case TrackSortFieldID, TrackSortFieldName, TrackSortFieldShortName, TrackSortFieldLength, TrackSortFieldPitlaneLength, TrackSortFieldNumSectors:
		return true
	}
	return false
}

func (e TrackSortField) String() string {
	return string(e)
}

func (e *TrackSortField) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TrackSortField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TrackSortField", str)
	}
	return nil
}

func (e TrackSortField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
