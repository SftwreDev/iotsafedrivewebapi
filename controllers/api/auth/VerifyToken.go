package auth

import (
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
	"time"
)

func VerifyTokenApi(w http.ResponseWriter, r *http.Request) {
	// Get userClaims data from request context
	userClaims, ok := utils.GetUserClaimsContext(r)
	if !ok {
		message := "User not found"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	email := userClaims.Email
	// Query database for user with provided email
	var appsUser []structs.Actor
	err := models.DB.Raw("SELECT * FROM apps_user WHERE email = ?", email).Scan(&appsUser).Error
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

	for _, user := range appsUser {
		// Generate access and refresh tokens
		userClaims := structs.UserClaims{
			ID:               user.ID,
			Email:            user.Email,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Role:             user.Role,
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
			Role:             user.Role,
			IsOnboardingDone: user.IsOnboardingDone,
			AccessToken:      signedAccessToken,
			RefreshToken:     signedRefreshToken,
		}
		// Send success response
		utils.SendSuccessResponse(http.StatusOK, "Access token still valid", response, w)
		return
	}
}
