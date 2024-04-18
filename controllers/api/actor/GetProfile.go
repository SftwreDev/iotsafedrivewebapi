package actor

import (
	"github.com/getsentry/sentry-go"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
)

func GetActorProfileApi(w http.ResponseWriter, r *http.Request) {

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

	// Query database for user with provided email
	var appsUser []structs.AppsUser
	err := models.DB.Raw(`
		SELECT 
			u.id,
			u.first_name,
			u.last_name,
			u.address,
			u.contact,
			u.email,
			u.device_id,
			u.profile_picture,
			u.date_joined,
			u.is_active,
			u.is_onboarding_done,
			u.is_staff,
			u.is_superuser,
			v.brand,
			v.model,
			v.year_model,
			v.plate_no
		FROM 
			apps_user AS u
		LEFT JOIN 
			apps_vehicle AS v
		ON 
			v.owner_id = u.id
		WHERE 
			u.email = ?`, userEmail).Scan(&appsUser).Error

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

		var appsTrustedContacts []structs.AppsTrustedContacts

		err := models.DB.Raw(`
			SELECT 
				id, name, contact, address
			FROM 
				apps_trustedcontacts
			WHERE 
				owner_id = ?`, user.ID).Scan(&appsTrustedContacts).Error

		if err != nil {
			sentry.CaptureException(err)
			// Return an error response if query fails
			utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		var trustedContacts []any

		for _, contact := range appsTrustedContacts {
			trustedContacts = append(trustedContacts, structs.AppsTrustedContacts{
				Name:    contact.Name,
				Address: contact.Address,
				Contact: contact.Contact,
			})
		}

		// Prepare response with user information and associated vehicles
		response := structs.AppsUser{
			ID:               user.ID,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Address:          user.Address,
			Contact:          user.Contact,
			Email:            user.Email,
			DeviceID:         user.DeviceID,
			ProfilePicture:   user.ProfilePicture,
			DateJoined:       user.DateJoined,
			IsActive:         user.IsActive,
			IsOnboardingDone: user.IsOnboardingDone,
			IsStaff:          user.IsStaff,
			IsSuperuser:      user.IsSuperuser,
			Brand:            user.Brand,
			Model:            user.Model,
			YearModel:        user.YearModel,
			PlateNo:          user.PlateNo, // Assign the vehicles to the response
			TrustedContacts:  trustedContacts,
		}

		// Send success response
		utils.SendSuccessResponse(http.StatusOK, "User's profile information", response, w)
	}

}
