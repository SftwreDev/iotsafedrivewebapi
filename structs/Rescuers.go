package structs

type Rescuers struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type SelectRescuer struct {
	ActivityHistoryID string `json:"activity_history_id"`
	RescuerID         string `json:"rescuer_id"`
	RespondersName    string `json:"responders_name"`
	Notes             string `json:"notes"`
}

type RescuerInformationDetails struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	Contact        string `json:"contact"`
	RespondersName string `json:"responders_name"`
	Notes          string `json:"notes"`
}

type AddRescuers struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}
