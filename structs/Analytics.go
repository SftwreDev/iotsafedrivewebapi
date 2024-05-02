package structs

type AnalyticsFilter struct {
	Filter string `json:"filter" `
	Value  string `json:"value"`
}

type ActivityHistoryStats struct {
	TotalActivity   int `json:"total_activity"`
	TotalSmsSent    int `json:"total_sms_sent"`
	TotalFalseAlarm int `json:"total_false_alarm"`
	TotalSmsNotSent int `json:"total_sms_not_sent"`
}

type AccidentRescueStats struct {
	TotalRespondedAccident int `json:"total_responded_accident"`
}

type SmsVsFalse struct {
	Timestamps      string `json:"timestamps"`
	TotalSmsSent    int    `json:"total_sms_sent"`
	TotalFalseAlarm int    `json:"total_false_alarm"`
	TotalSmsNotSent int    `json:"total_sms_not_sent"`
}

type AccidentRescuer struct {
	Timestamps           string `json:"timestamps"`
	TotalRescuedAccident int    `json:"total_rescued_accident"`
}

type OnGoingVsClosedVsPending struct {
	PendingPercentage float64 `json:"pending_percentage"`
	ClosedPercentage  float64 `json:"closed_percentage"`
	OnGoingPercentage float64 `json:"on_going_percentage"`
}
