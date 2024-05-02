package actor

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"io"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
	"time"
)

func AddAccountApi(w http.ResponseWriter, r *http.Request) {

	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	var payload structs.NewAccount
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &payload)

	// Validate input using validator package
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		sentry.CaptureException(err)
		// Return a validation error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Get timestamp today
	dateJoined := time.Now()

	// et role constants value
	role := utils.GetRole(payload.Role)

	// Hash the password securely
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	err = models.DB.Exec(`
		INSERT INTO apps_user(
		                      email,
		                      first_name,
		                      last_name,
		                      address,
		                      contact,
		                      role,
		                      date_joined,
		                      password,
		                      username,
		                      device_id,
		                      profile_picture,
		                      is_onboarding_done
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		payload.Email,
		payload.FirstName,
		payload.LastName,
		payload.Address,
		payload.Contact,
		role,
		dateJoined,
		hashedPassword,
		"",
		"",
		"",
		true,
	).Error

	if err != nil {
		sentry.CaptureException(err)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Return response message
	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully added new account", response, w)
	return

}

func IsPasswordChangedApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get userClaims data from request context
	userClaims, ok := utils.GetUserClaimsContext(r)
	if !ok {
		message := "User not found "

		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	// Declare userEmail from userClaims
	userEmail := userClaims.Email

	var isPasswordChanged bool

	execQuery := models.DB.Raw(`
		SELECT is_password_changed FROM apps_user
		WHERE email = ?
	`, userEmail).Scan(&isPasswordChanged).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	if isPasswordChanged {
		utils.SendSuccessResponse(http.StatusOK, "Password already changed", isPasswordChanged, w)
		return
	} else {
		utils.SendSuccessResponse(http.StatusOK, "Password not yet change", isPasswordChanged, w)
		return
	}

}
