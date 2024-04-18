package structs

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

type SendSMS struct {
	IsFalseAlarm bool   `json:"is_false_alarm"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
}
