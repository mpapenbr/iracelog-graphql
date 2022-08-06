package analysis

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DbAnalysis struct {
	ID int `json:"id"`
	// TODO: Cars
	CarInfo []struct {
		Name     string `json:"name"`
		CarNum   string `json:"carNum"`
		CarClass string `json:"carClass"`
		Drivers  []struct {
			SeatTime []struct {
				EnterCarTime float64 // unit: sessionTime
				LeaveCarTime float64 // unit: sessionTime
			}
			DriverName string `json:"driverName"`
		}
	}
	CarLaps []struct {
		CarNum string
		Laps   []struct {
			LapNo   int
			LapTime float64
		}
	}
	RaceOrder []string // last computed race order (contain carNums)
}

func GetAnalysisForEvent(pool *pgxpool.Pool, eventId int) (*DbAnalysis, error) {

	var data DbAnalysis
	pool.QueryRow(context.Background(), "select id,data from analysis where event_id=$1", eventId).Scan(&data.ID, &data)
	return &data, nil

}
