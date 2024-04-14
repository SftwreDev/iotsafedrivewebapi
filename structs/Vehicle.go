package structs

type Vehicle struct {
	ID uint `gorm:"primary_key;autoIncrement"`

	Brand     string `json:"brand" validate:"required"`
	Model     string `json:"model" validate:"required"`
	YearModel string `json:"year_model" validate:"required"`
	PlateNo   string `json:"plate_no" validate:"required"`
}

type AllVehicle struct {
	OwnerID   uint   `json:"owner_id"`
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	YearModel string `json:"year_model"`
	PlateNo   string `json:"plate_no"`
	Owner     string `json:"owner"`
}
