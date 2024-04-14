package auth

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
	"reflect"
)

func ResetPasswordApi(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	var input structs.ResetPassword
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Validate input using validator package
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		// Return a validation error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Query database for user with provided email
	var appsUser structs.Actor
	err = models.DB.Raw("SELECT email FROM apps_user WHERE email = ?", input.Email).Scan(&appsUser).Error
	if err != nil {
		// Return an error response if query fails
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Check if appUser is empty
	if reflect.DeepEqual(appsUser, structs.Actor{}) {
		message := "User not found for the given email"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	// Hash the password securely
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	result := models.DB.Exec(`
		UPDATE apps_user
		SET password = ?
		WHERE email = ?
	`, hashedPassword, input.Email).Error

	if result != nil {
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error(), w)
		return
	}

	// Return response message
	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully reset your password", response, w)
	return
}
