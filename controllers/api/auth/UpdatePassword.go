package auth

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"io"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
)

func UpdatePasswordApi(w http.ResponseWriter, r *http.Request) {

	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get userClaims data from request context
	userClaims, ok := utils.GetUserClaimsContext(r)
	if !ok {
		message := "User not found "
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	// Declare userID from userClaims
	userID := userClaims.ID

	var input structs.UpdatePassword
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

	// Hash the password securely
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	result := models.DB.Exec(`
		UPDATE apps_user
		SET password = ?
		WHERE id = ?
	`, hashedPassword, userID).Error

	if result != nil {
		sentry.CaptureException(result)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error(), w)
		return
	}

	// Return response message
	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully updated your password", response, w)
	return
}

func UpdateTemporaryPassword(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	var input structs.UpdateTempPassword
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

	// Get userClaims data from request context
	userClaims, ok := utils.GetUserClaimsContext(r)
	if !ok {
		message := "User not found "
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	// Declare userID from userClaims
	userID := userClaims.ID

	// Query database for user with provided email
	var appsUser structs.Actor
	err = models.DB.Raw("SELECT * FROM apps_user WHERE id = ?", userID).Scan(&appsUser).Error
	if err != nil {
		sentry.CaptureException(err)
		// Return an error response if query fails
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	validPassword := utils.CheckPasswordHash(input.CurrentPassword, appsUser.Password)
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	if validPassword {
		err = models.DB.Exec(`
				UPDATE apps_user
				SET password = ?,
					is_password_changed = ? 
				WHERE id = ?
			`, hashedPassword, true, userID).Error
		if err != nil {
			sentry.CaptureException(err)
			// Return an error response if query fails
			utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		// Send success response
		var response []interface{}
		utils.SendSuccessResponse(http.StatusCreated, "Successfully updated your password.", response, w)
		return
	} else {
		// Return an error response if query fails
		utils.SendErrorResponse(http.StatusBadRequest, "Invalid current password", w)
		return
	}

}
