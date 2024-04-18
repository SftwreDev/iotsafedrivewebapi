package trusted_contacts

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

func AddTrustedContactsApi(w http.ResponseWriter, r *http.Request) {
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

	var input []structs.TrustedContacts
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Validate input using validator package
	validate := validator.New()

	for _, contact := range input {
		err := validate.Struct(contact)
		if err != nil {
			// Return a validation error response
			utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		result := models.DB.Exec(`
			INSERT INTO apps_trustedcontacts(name, address, contact,  owner_id) 
			VALUES (?, ?, ?, ?)
		`, contact.Name, contact.Address, contact.Contact, userID).Error

		if result != nil {
			sentry.CaptureException(result)
			// Return an error response
			utils.SendErrorResponse(http.StatusBadRequest, result.Error(), w)
			return
		}
	}

	result := models.DB.Exec(`
		UPDATE apps_user
		SET is_onboarding_done = true
		WHERE id = ?
	`, userID).Error

	if result != nil {
		sentry.CaptureException(result)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error(), w)
		return
	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully added your trusted contacts", response, w)
	return
}

func ListTrustedContactsApi(w http.ResponseWriter, r *http.Request) {
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
	var trustedContacts []structs.TrustedContacts
	err := models.DB.Raw(`
		SELECT 
			id, name, contact, address
		FROM 
			apps_trustedcontacts
		WHERE 
			owner_id = ?`, userID).Scan(&trustedContacts).Error

	if err != nil {
		sentry.CaptureException(err)
		// Return an error response if query fails
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Check if user with provided email exists
	if len(trustedContacts) == 0 {
		message := "Trusted contacts not found"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	// Send success response
	utils.SendSuccessResponse(http.StatusOK, "User's trusted contacts", trustedContacts, w)
	return
}

func UpdateTrustedContactsApi(w http.ResponseWriter, r *http.Request) {

	var input []structs.TrustedContacts
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Validate input using validator package
	validate := validator.New()

	for _, contact := range input {
		err := validate.Struct(contact)
		if err != nil {
			sentry.CaptureException(err)
			// Return a validation error response
			utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		result := models.DB.Exec(`
			UPDATE apps_trustedcontacts
			SET 
			    name = ?,
			    address = ?,
			    contact = ?
			WHERE
				id = ?
		`, contact.Name, contact.Address, contact.Contact, contact.ID)

		if result.Error != nil {
			sentry.CaptureException(result.Error)
			// Return an error response
			utils.SendErrorResponse(http.StatusBadRequest, result.Error.Error(), w)
			return
		}

	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully updated your trusted contacts", response, w)
	return
}

func GetAllTrustedContactsApi(w http.ResponseWriter, r *http.Request) {

	var trustedContacts []structs.AllTrustedContacts

	err := models.DB.Raw(
		`
				SELECT
					c.owner_id, c.name, c.address, c.contact,
					CONCAT(u.first_name, ' ', u.last_name) AS owner
				FROM apps_trustedcontacts as c
				INNER JOIN apps_user as u ON u.id = c.owner_id;
			`,
	).Scan(&trustedContacts).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved your vehicle info", trustedContacts, w)
	return
}
