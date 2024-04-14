package models

type Rescuers struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

// TableName specifies the table name for the model.
func (Rescuers) TableName() string {
	return "apps_rescueteamcontacts"
}
