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
	DeviceID     string    `json:"device_id"`
}

type ForwardedAccidents struct {
	ID                 string    `json:"id"`
	Notes              string    `json:"notes"`
	Status             string    `json:"status"`
	ForwardedBy        string    `json:"forwarded_by"`
	Location           string    `json:"location"`
	Victim             string    `json:"victim"`
	ActivityHistoryID  string    `json:"activity_history_id"`
	ForwardedOn        time.Time `json:"forwarded_on"`
	AccidentOccurredOn time.Time `json:"accident_occurred_on"`
}
