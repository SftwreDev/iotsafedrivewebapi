package auth

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"io"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
	"time"
)

// SignInApi handles the sign-in functionality.
// It expects a POST request with a JSON body containing email and password.
// If the credentials are valid, it generates access and refresh tokens.
// It returns a JSON response with the user information and tokens.
func SignInApi(w http.ResponseWriter, r *http.Request) {

	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Read and parse request body into SignInInput struct
	var input structs.SignInInput
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Validate input using validator package
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		sentry.CaptureException(err)
		// Return a validation error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Query database for user with provided email
	var appsUser []structs.Actor
	err = models.DB.Raw("SELECT * FROM apps_user WHERE email = ?", input.Email).Scan(&appsUser).Error
	if err != nil {
		sentry.CaptureException(err)
		// Return an error response if query fails
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Check if user with provided email exists
	if len(appsUser) == 0 {
		message := "User not found for the given email"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	// Check if the provided password matches the hashed password in the database
	var validPassword bool
	for _, user := range appsUser {
		validPassword = utils.CheckPasswordHash(input.Password, user.Password)

		if validPassword {

			// Generate access and refresh tokens
			userClaims := structs.UserClaims{
				ID:               user.ID,
				Email:            user.Email,
				FirstName:        user.FirstName,
				LastName:         user.LastName,
				IsStaff:          user.IsStaff,
				IsActive:         user.IsActive,
				IsSuperuser:      user.IsSuperuser,
				IsOnboardingDone: user.IsOnboardingDone,
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(time.Minute * 180).Unix(),
				},
			}
			signedAccessToken, signedRefreshToken, err := utils.GenerateAccessAndRefreshTokens(userClaims)
			if err != nil {
				sentry.CaptureException(err)
				// Return an error response if token generation fails
				utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
				return
			}

			// Prepare response with user information and tokens
			response := structs.SignInOutput{
				ID:               user.ID,
				Email:            user.Email,
				FirstName:        user.FirstName,
				LastName:         user.LastName,
				ProfilePicture:   user.ProfilePicture,
				IsStaff:          user.IsStaff,
				IsSuperuser:      user.IsSuperuser,
				IsOnboardingDone: user.IsOnboardingDone,
				AccessToken:      signedAccessToken,
				RefreshToken:     signedRefreshToken,
			}
			// Send success response
			utils.SendSuccessResponse(http.StatusOK, "Successfully logged in", response, w)
			return
		} else {
			// If password does not match, return incorrect password error response
			message := "Invalid email or password"
			utils.SendErrorResponse(http.StatusBadRequest, message, w)
		}
	}

}
