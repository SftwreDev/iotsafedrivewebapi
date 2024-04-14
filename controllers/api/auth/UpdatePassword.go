package auth

import (
	"encoding/json"
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
		// Return a validation error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
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
		WHERE id = ?
	`, hashedPassword, userID).Error

	if result != nil {
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error(), w)
		return
	}

	// Return response message
	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully updated your password", response, w)
	return
}
