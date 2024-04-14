package actor

import (
	"fmt"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"log"
	"net/http"
)

func UpdateActorProfileApi(w http.ResponseWriter, r *http.Request) {

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

	// Parse form values from the request
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	address := r.FormValue("address")
	contact := r.FormValue("contact")
	deviceID := r.FormValue("device_id")

	// Upload profile picture to TempFile from FormData
	file, err := utils.UploadFileFromFormData("profile_picture", r)

	// Declare string variables for cloudinary secureURL and publicID
	var secureURL, publicID string

	if err != nil {
		// Log the error and return it to the caller
		log.Printf("File not found: %v", err.Error())

	} else {
		// Upload profile picture to Cloudinary
		secureURL, publicID, err = utils.UploadPublicFileToCloudinary(file)

		if err != nil {
			// Log the error and return it to the caller
			log.Printf("Failed to upload file to Cloudinary: %v", err)
			// Return an error response
			utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}
	}

	result := models.DB.Exec(`
        UPDATE apps_user
        SET 
            email = ?,
            first_name = ?,
            last_name = ?,
            address = ?,
            contact = ?,
            device_id = ?,
            profile_picture = CASE WHEN ? <> '' THEN ? ELSE profile_picture END
        WHERE
            id = ?
    `, email, firstName, lastName, address, contact, deviceID, secureURL, secureURL, userID)

	// Check if the operation was successful
	if result.Error != nil {
		// If there is an error, delete the uploaded profile picture from Cloudinary
		deleteFile, _ := utils.DeleteFileFromCloudinary(publicID)
		fmt.Println(deleteFile)

		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error.Error(), w)
		return
	}

	response := structs.UpdateActor{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Address:        address,
		Contact:        contact,
		DeviceID:       deviceID,
		ProfilePicture: secureURL,
	}

	// Create a success response with the user data
	utils.SendSuccessResponse(http.StatusCreated, "Successfully update your profile", response, w)

}
