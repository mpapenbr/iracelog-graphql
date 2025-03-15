package mytypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

//nolint:tagliatelle // json is that way
type (
	// we need this extra type to be usable with bob generated code
	EventSessionSlice []EventSession
	EventSession      struct {
		Num         int    `json:"num"`
		Name        string `json:"name"`
		Type        int    `json:"type"`
		SessionTime int    `json:"session_time"`
		Laps        int    `json:"laps"`
	}
	// we need this extra type to be usable with bob generated code
	SectorSlice []Sector
	Sector      struct {
		Num      int     `json:"num"`
		StartPct float64 `json:"start_pct"`
	}
)

func (h *EventSessionSlice) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte")
	}

	return json.Unmarshal(bytes, &h)
}

func (h EventSessionSlice) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *SectorSlice) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte")
	}

	return json.Unmarshal(bytes, &h)
}

func (h SectorSlice) Value() (driver.Value, error) {
	return json.Marshal(h)
}
