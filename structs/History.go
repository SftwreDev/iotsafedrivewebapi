package structs

import "time"

type ActivityHistory struct {
	ID           string    `json:"id"`
	TimeStamps   time.Time `json:"time_stamps"`
	Location     string    `json:"location"`
	Latitude     string    `json:"latitude"`
	Longitude    string    `json:"longitude"`
	Status       string    `json:"status"`
	StatusReport string    `json:"status_report"`
	Owner        string    `json:"owner"`
}
