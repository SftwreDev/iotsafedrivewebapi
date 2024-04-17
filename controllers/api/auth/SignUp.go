package auth

import (
	"fmt"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
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

	// Hash the password securely
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Upload profile picture to TempFile from FormData
	file, err := utils.UploadFileFromFormData("profile_picture", r)

	if err != nil {
		// Log the error and return it to the caller
		log.Printf("Failed to upload file to Cloudinary: %v", err)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Upload profile picture to Cloudinary
	secureURL, publicID, err := utils.UploadPublicFileToCloudinary(file)

	if err != nil {
		// Log the error and return it to the caller
		log.Printf("Failed to upload file to Cloudinary: %v", err)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Create a new user object with form values
	newUser := &models.AppsUser{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Address:        address,
		Contact:        contact,
		DeviceID:       deviceID,
		DateJoined:     time.Now(),
		ProfilePicture: secureURL,
		Password:       hashedPassword,
	}

	// Save user data to the database
	result := models.DB.Create(newUser)
	if result.Error != nil {

		// If there is an error, delete the uploaded profile picture from Cloudinary
		deleteFile, _ := utils.DeleteFileFromCloudinary(publicID)
		fmt.Println(deleteFile)

		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error.Error(), w)
		return
	}

	// Convert the user object to response format
	response := structs.AppsUser{
		ID:               newUser.ID,
		FirstName:        newUser.FirstName,
		LastName:         newUser.LastName,
		Address:          newUser.Address,
		Contact:          newUser.Contact,
		Email:            newUser.Email,
		DeviceID:         newUser.DeviceID,
		ProfilePicture:   newUser.ProfilePicture,
		IsActive:         newUser.IsActive,
		IsOnboardingDone: newUser.IsOnboardingDone,
		IsStaff:          newUser.IsStaff,
		IsSuperuser:      newUser.IsSuperuser,
	}

	// Create a success response with the user data
	utils.SendSuccessResponse(http.StatusCreated, "Success", response, w)
}
