package structs

type Actor struct {
	ID uint `json:"id"`

	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Address        string `json:"address"`
	Contact        string `json:"contact"`
	Email          string `json:"email"`
	DeviceID       string `json:"device_id"`
	ProfilePicture string `json:"profile_picture"`
	DateJoined     string `json:"date_joined"`
	Role           string `json:"role"`
	Password       []byte `json:"_"`

	IsOnboardingDone bool `json:"is_onboarding_done"`
}

type UpdateActor struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Address        string `json:"address"`
	Contact        string `json:"contact"`
	Email          string `json:"email"`
	DeviceID       string `json:"device_id"`
	ProfilePicture string `json:"profile_picture"`
}

type AllUsers struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address    string `json:"address"`
	Contact    string `json:"contact"`
	Email      string `json:"email"`
	DeviceID   string `json:"device_id"`
	Role       string `json:"role"`
	DateJoined string `json:"date_joined"`
}

func (u *AllUsers) Modify() {
	if u.Role == "super_admin" {
		u.Role = "Super Admin"
	}
}
