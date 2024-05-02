package auth

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"iotsafedriveapi/models"
	"iotsafedriveapi/utils"
	"log"
	"net/http"
	"time"
)

// SignUpApi is the API endpoint for user sign up
func SignUpApi(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse form values from the request
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	address := r.FormValue("address")
	contact := r.FormValue("contact")
	deviceID := r.FormValue("device_id")
	role := "user"

	// Get timestamps today
	dateJoined := time.Now()

	// Hash the password securely
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Upload profile picture to TempFile from FormData
	file, err := utils.UploadFileFromFormData("profile_picture", r)

	if err != nil {
		sentry.CaptureException(err)
		// Log the error and return it to the caller
		log.Printf("Failed to upload file to Cloudinary: %v", err)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Upload profile picture to Cloudinary
	secureURL, publicID, err := utils.UploadPublicFileToCloudinary(file)

	if err != nil {
		sentry.CaptureException(err)
		// Log the error and return it to the caller
		log.Printf("Failed to upload file to Cloudinary: %v", err)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Create a new user object with form values
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
		                      is_onboarding_done,
		                      is_password_changed
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true)
	`,
		email,
		firstName,
		lastName,
		address,
		contact,
		role,
		dateJoined,
		hashedPassword,
		"",
		deviceID,
		secureURL,
		false,
	).Error

	if err != nil {
		sentry.CaptureException(err)
		// If there is an error, delete the uploaded profile picture from Cloudinary
		deleteFile, _ := utils.DeleteFileFromCloudinary(publicID)
		fmt.Println(deleteFile)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Return response message
	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully created your account", response, w)
	return
}
