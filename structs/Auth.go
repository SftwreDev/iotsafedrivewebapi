package structs

import (
	"github.com/golang-jwt/jwt"
)

type SignUpInput struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Contact   string `json:"contact" validate:"required"`
	Email     string `json:"email" validate:"required"`
	DeviceID  string `json:"device_id" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInOutput struct {
	ID               uint   `json:"id"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	ProfilePicture   string `json:"profile_picture"`
	Role             string `json:"role"`
	IsOnboardingDone bool   `json:"is_onboarding_done"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
}

type UserClaims struct {
	ID               uint   `json:"id"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Role             string `json:"role"`
	IsOnboardingDone bool   `json:"is_onboarding_done"`
	jwt.StandardClaims
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ResetPassword struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdatePassword struct {
	Password string `json:"password" validate:"required"`
}

func (userClaims UserClaims) UserClaimsPointer() UserClaims {
	return userClaims
}

func (userClaims *UserClaims) SetUserClaimsPointer(email string, firstName string, lastName string) {
	userClaims.Email = email
	userClaims.FirstName = firstName
	userClaims.LastName = lastName
}

type UpdateTempPassword struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}
