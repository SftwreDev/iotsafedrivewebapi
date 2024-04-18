package accident_alert

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"io"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
	"time"
)

func GetAllAccidentAlertApi(w http.ResponseWriter, r *http.Request) {

	var accidentAlert []structs.AccidentAlert

	err := models.DB.Raw(
		`
				SELECT
					a.latitude, a.longitude, a.is_active, a.device_id,
					CONCAT(u.first_name, ' ', u.last_name) AS owner
				FROM apps_accidentalert as a
				INNER JOIN apps_user as u ON u.device_id = a.device_id;
			`,
	).Scan(&accidentAlert).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Check if user with provided email exists
	if len(accidentAlert) == 0 {
		message := "Accident alert not found"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	for _, alert := range accidentAlert {
		response := structs.AccidentAlertOutput{
			Latitude:  alert.Latitude,
			Longitude: alert.Longitude,
			DeviceID:  alert.DeviceID,
			IsActive:  alert.IsActive,
			User:      alert.Owner,
			Location:  utils.GetLocation(alert.Latitude, alert.Longitude),
		}

		utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved accident alerts", response, w)
		return
	}

}

func AccidentDetectedApi(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Read and parse request body into SignInInput struct
	var payload structs.IoTAlert
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

	execQuery := models.DB.Exec(`
		INSERT INTO apps_accidentalert(
			latitude, 
			longitude, 
			device_id,
			is_active
		)
		VALUES (?, ?, ?, ?)
	`, payload.Latitude, payload.Longitude, payload.DeviceID, true).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	var empty []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "successfully send alert", empty, w)
	return
}

func GetLatestAccidentAlertApi(w http.ResponseWriter, r *http.Request) {

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

	var deviceID string

	execQuery := models.DB.Raw(`
		SELECT device_id from apps_user
		WHERE id = ?
	`, userClaims.ID).Scan(&deviceID).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	var alert []structs.IoTAlert

	err := models.DB.Raw(
		`
				SELECT
					latitude, longitude, is_active, device_id
				FROM apps_accidentalert
				WHERE device_id = ?
				LIMIT 1
			`, deviceID).Scan(&alert).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	if len(alert) == 0 {
		message := "No pending accident as of the moment"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusOK, message, w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved latest activity history", alert, w)
	return
}

func SendSMSApi(w http.ResponseWriter, r *http.Request) {
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

	// Declare userEmail from userClaims
	userEmail := userClaims.Email

	// Read and parse request body into SignInInput struct
	var payload structs.SendSMS
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

	var userDeviceID string

	execQuery := models.DB.Raw(`
		SELECT device_id from apps_user where email = ?
	`, userEmail).Scan(&userDeviceID).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	var alert structs.IoTAlert

	execQuery = models.DB.Raw(
		`
				SELECT
					latitude, longitude, is_active, device_id
				FROM apps_accidentalert
				WHERE  device_id = ?
				LIMIT 1
			`, userDeviceID).Scan(&alert).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	// Query database for user with provided email
	var trustedContacts []structs.TrustedContacts
	execQuery = models.DB.Raw(`
		SELECT 
			id, name, contact, address
		FROM 
			apps_trustedcontacts
		WHERE 
			owner_id = ?`, userClaims.ID).Scan(&trustedContacts).Error

	if execQuery != nil {
		// Return an error response if query fails
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	isFalseAlarm := payload.IsFalseAlarm

	var latitude string
	var longitude string

	if alert.Latitude == "" || alert.Longitude == "" {
		latitude = payload.Latitude
		longitude = payload.Longitude
	} else {
		latitude = alert.Latitude
		longitude = alert.Longitude
	}

	timestamps := time.Now()
	location := utils.GetLocation(latitude, longitude)
	userID := userClaims.ID

	if isFalseAlarm {

		execQuery = models.DB.Exec(`
			INSERT INTO apps_activityhistory(
		 		 timestamps, 
				 location, 
				 latitude, 
				 longitude, 
				 status, 
				 user_id, 
				 message, 
				 status_report
			) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, timestamps, location, latitude, longitude, "False Alarm", userID, "", "closed").Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}

		execQuery = models.DB.Exec(`
			DELETE FROM apps_accidentalert
			WHERE device_id = ?
		`, userDeviceID).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}
	} else {
		message := fmt.Sprintf(
			`URGENT: Accident! This is %s %s. I'm in an accident', need ambulance. Pls send help ASAP! Location: %s`,
			userClaims.FirstName, userClaims.LastName, location)

		var status string
		for _, contact := range trustedContacts {
			fmt.Printf("\n Sending sms now to %s", contact.Contact)

			sendSMS := utils.SendSMS(
				contact.Contact,
				message,
			)

			if sendSMS != nil {
				message = sendSMS.Error()
				status = "SMS Not Sent"
			} else {
				status = "SMS Sent"
			}
		}
		execQuery = models.DB.Exec(`
			INSERT INTO apps_activityhistory(
				 timestamps, 
				 location, 
				 latitude, 
				 longitude, 
				 status, 
				 user_id, 
				 message, 
				 status_report
			) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, timestamps, location, latitude, longitude, status, userID, message, "pending").Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}

		execQuery = models.DB.Exec(`
			DELETE FROM apps_accidentalert
			WHERE device_id = ?
		`, userDeviceID).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}

		var response []interface{}
		utils.SendSuccessResponse(http.StatusOK, status, response, w)
		return
	}
}
