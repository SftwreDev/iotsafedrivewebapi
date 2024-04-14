package models

import (
	"time"
)

type AppsUser struct {
	ID uint `json:"id"`

	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Address        string    `json:"address"`
	Contact        string    `json:"contact"`
	Email          string    `json:"email"`
	DeviceID       string    `json:"device_id"`
	ProfilePicture string    `json:"profile_picture"`
	DateJoined     time.Time `json:"date_joined"`
	Password       []byte    `json:"_"`

	IsActive         bool `json:"is_active"`
	IsOnboardingDone bool `json:"is_onboarding_done"`
	IsStaff          bool `json:"is_staff"`
	IsSuperuser      bool `json:"is_superuser"`
}

// TableName specifies the table name for the model.
func (AppsUser) TableName() string {
	return "apps_user"
}
