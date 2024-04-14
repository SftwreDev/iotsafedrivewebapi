package structs

type AppsUser struct {
	ID uint `json:"id"`

	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Address        string `json:"address"`
	Contact        string `json:"contact"`
	Email          string `json:"email"`
	DeviceID       string `json:"device_id"`
	ProfilePicture string `json:"profile_picture"`
	DateJoined     string `json:"date_joined"`

	IsActive         bool `json:"is_active"`
	IsOnboardingDone bool `json:"is_onboarding_done"`
	IsStaff          bool `json:"is_staff"`
	IsSuperuser      bool `json:"is_superuser"`

	Brand           string        `json:"brand"`
	Model           string        `json:"model"`
	YearModel       string        `json:"year_model"`
	PlateNo         string        `json:"plate_no"`
	TrustedContacts []interface{} `json:"trusted_contacts"`
}

type AppsTrustedContacts struct {
	ID uint `gorm:"primary_key;autoIncrement"`

	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func (appsUser *AppsUser) AddItem(item AppsTrustedContacts) {
	appsUser.TrustedContacts = append(appsUser.TrustedContacts, item)
}
