package structs

import "time"

type AccidentAlert struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	DeviceID  string `json:"device_id"`
	IsActive  bool   `json:"is_active"`
	Owner     string `json:"owner"`
}

type AccidentAlertOutput struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	DeviceID  string `json:"device_id"`
	IsActive  bool   `json:"is_active"`

	User     string `json:"user"`
	Location string `json:"location"`
}

type IoTAlert struct {
	DeviceID  string `json:"device_id" validate:"required"`
	Latitude  string `json:"latitude" validate:"required"`
	Longitude string `json:"longitude" validate:"required"`
}

type SendSMSStructs struct {
	IsFalseAlarm bool   `json:"is_false_alarm"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
}

type ForwardAccident struct {
	RescuerID         string `json:"rescuer_id"`
	ActivityHistoryID string `json:"activity_history_id"`
	Notes             string `json:"notes"`
	Status            string `json:"status"`
	ForwardedBy       string `json:"forwarded_by"`
}

type RejectedNotifications struct {
	Reason string `json:"reason"`
}

type AcceptedAccidents struct {
	ID         int       `json:"id"`
	Notes      string    `json:"notes"`
	Status     string    `json:"status"`
	AcceptedBy string    `json:"accepted_by"`
	Rescuer    string    `json:"rescuer"`
	Patient    string    `json:"patient"`
	Timestamps time.Time `json:"timestamps"`
}

type RejectedAccidents struct {
	ID         int       `json:"id"`
	Notes      string    `json:"notes"`
	Status     string    `json:"status"`
	RejectedBy string    `json:"rejected_by"`
	Rescuer    string    `json:"rescuer"`
	Patient    string    `json:"patient"`
	Timestamps time.Time `json:"timestamps"`
}
