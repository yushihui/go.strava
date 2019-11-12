package strava

import "time"

type Split struct {
	Distance            float64   `json:"distance"`
	ElapsedTime         int       `json:"elapsed_time"`
	ElevationDifference float64   `json:"elevation_difference"`
	MovingTime          int       `json:"moving_time"`
	Split               int       `json:"split"`
	AvgHr               float64   `json:"average_heartrate"`
	AvgSpeed            float64   `json:"average_speed"`
	IsRace              bool      `json:"is_race"`
	Temperature         int       `json:"temperature"`
	WindSpeed           int       `json:"wind_speed"`
	Humidity            int       `json:"humidity"`
	Invalid             bool      `json:"invalid"`
	StartDate           time.Time `json:"start_date"`
	ActivityId          int64     `json:"activity_id"`
}
