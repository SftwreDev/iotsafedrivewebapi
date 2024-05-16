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
	"strconv"
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

	var deviceID string
	execQuery := models.DB.Raw(`SELECT device_id FROM apps_accidentalert WHERE device_id = ?`, payload.DeviceID).Scan(&deviceID).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	if deviceID != "" {
		execQuery = models.DB.Exec(`DELETE FROM apps_accidentalert WHERE device_id = ?`, payload.DeviceID).Error
		if execQuery != nil {
			sentry.CaptureException(execQuery)
			// Return an error response
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}
	}
	execQuery = models.DB.Exec(`
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
	var payload structs.SendSMSStructs
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &payload)

	fmt.Println(payload)

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

	latitude := payload.Lat
	longitude := payload.Lng

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
		`, timestamps, location, latitude, longitude, "False Alarm", userID, "None", "closed").Error

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

func ForwardAccidentApi(w http.ResponseWriter, r *http.Request) {

	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Read and parse request body into SignInInput struct
	var payload structs.ForwardAccident
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
	INSERT INTO apps_forwarded_accidents(rescuer_id, notes, status, forwarded_by, activity_history_id)
	VALUES(?, ?, ?, ?, ?)
	`, payload.RescuerID, payload.Notes, payload.Status, payload.ForwardedBy, payload.ActivityHistoryID).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	execQuery = models.DB.Exec(`
		UPDATE apps_activityhistory 
		SET status_report = 'forwarded'
		WHERE id = ?
	`, payload.ActivityHistoryID).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	go utils.PerformPendingAccidents(payload.ActivityHistoryID)

	var response []interface{}
	utils.SendSuccessResponse(http.StatusOK, "Successfully forwarded accident", response, w)
	return
}

func CheckIfAccidentIsForwarded(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Parse query parameters from the request URL
	queryValues := r.URL.Query()

	activityId := queryValues.Get("id")

	var activityHistoryID string

	execQuery := models.DB.Raw(`
			SELECT activity_history_id FROM apps_forwarded_accidents
			WHERE activity_history_id = ? AND status = 'pending'
			ORDER BY timestamps DESC
			LIMIT 1
		`, activityId).Scan(&activityHistoryID).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	if activityHistoryID == "" {
		utils.SendErrorResponse(http.StatusNotFound, "Not yet forwarded", w)
		return
	} else {
		var response []interface{}
		utils.SendSuccessResponse(http.StatusOK, "Accident already forwarded", response, w)
		return
	}

}

func ForwardedAccidentsActions(w http.ResponseWriter, r *http.Request) {
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
	rescuer := fmt.Sprintf(`%s %s`, userClaims.FirstName, userClaims.LastName)
	userID := userClaims.ID

	// Parse query parameters from the request URL
	queryValues := r.URL.Query()

	action := queryValues.Get("action")
	activityId := queryValues.Get("activity_id")

	// Read and parse request body into SignInInput struct
	var payload structs.RejectedNotifications
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

	if action == "rejected" {
		execQuery := models.DB.Exec(`
			DELETE FROM apps_forwarded_accidents WHERE activity_history_id = ?
		`, activityId).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}

		execQuery = models.DB.Exec(`
			INSERT INTO apps_rejected_accidents(rescuer_id, activity_history_id, notes, status, rejected_by) 
			VALUES (?, ?, ?, ?, ?)`, strconv.Itoa(int(userID)), activityId, payload.Reason, "rejected", rescuer).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}

		execQuery = models.DB.Exec(`
			UPDATE apps_activityhistory
			SET status_report = ?
			WHERE id = ?
		`, action, activityId).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}
	} else {
		execQuery := models.DB.Exec(`
			INSERT INTO apps_accepted_accidents(rescuer_id, activity_history_id, notes, status, accepted_by) 
			VALUES (?, ?, ?, ?, ?)`, strconv.Itoa(int(userID)), activityId, payload.Reason, "accepted", rescuer).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}

		execQuery = models.DB.Exec(`
			UPDATE apps_activityhistory
			SET status_report = ?
			WHERE id = ?
		`, action, activityId).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
			utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
			return
		}
	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusOK, "", response, w)
	return
}

func AcceptedAccidents(w http.ResponseWriter, r *http.Request) {
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

	userID := userClaims.ID

	var role string

	execQuery := models.DB.Raw(`SELECT role FROM apps_user WHERE id = ?`, userID).Scan(&role).Error
	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}
	var acceptedAccidents []structs.AcceptedAccidents

	var sqlQuery string

	baseQuery := `
		SELECT
			ac.id,
			ac.notes,
			ac.status,
			ac.accepted_by,
			ac.timestamps,
			CONCAT(u.first_name, ' ', u.last_name) AS rescuer,
			CONCAT(us.first_name, ' ', us.last_name) AS patient
		FROM
			apps_accepted_accidents AS ac
				INNER JOIN
			apps_user AS u ON u.id = ac.rescuer_id
				INNER JOIN
			apps_activityhistory AS ah ON ah.id = ac.activity_history_id
				INNER JOIN
			apps_user AS us ON us.id = ah.user_id
	`

	whereClause := fmt.Sprintf("WHERE ac.rescuer_id = %d", userID)
	orderByClause := ` ORDER BY ac.timestamps DESC`

	switch role {
	case "rescuer":
		sqlQuery = baseQuery + whereClause + orderByClause

	default:
		sqlQuery = baseQuery + orderByClause
	}

	execQuery = models.DB.Raw(sqlQuery).Scan(&acceptedAccidents).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "", acceptedAccidents, w)
	return
}

func RejectedAccidents(w http.ResponseWriter, r *http.Request) {
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

	userID := userClaims.ID

	var role string

	execQuery := models.DB.Raw(`SELECT role FROM apps_user WHERE id = ?`, userID).Scan(&role).Error
	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	var rejectedAccidents []structs.RejectedAccidents

	var sqlQuery string

	baseQuery := `
		SELECT
			rj.id,
			rj.notes,
			rj.status,
			rj.rejected_by,
			rj.timestamps,
			CONCAT(u.first_name, ' ', u.last_name) AS rescuer,
			CONCAT(us.first_name, ' ', us.last_name) AS patient
		FROM
			apps_rejected_accidents AS rj
				INNER JOIN
			apps_user AS u ON u.id = rj.rescuer_id
				INNER JOIN
			apps_activityhistory AS ah ON ah.id = rj.activity_history_id
				INNER JOIN
			apps_user AS us ON us.id = ah.user_id
	`

	whereClause := fmt.Sprintf("WHERE rj.rescuer_id = %d", userID)
	orderByClause := ` ORDER BY rj.timestamps DESC`

	switch role {
	case "rescuer":
		sqlQuery = baseQuery + whereClause + orderByClause

	default:
		sqlQuery = baseQuery + orderByClause
	}

	execQuery = models.DB.Raw(sqlQuery).Scan(&rejectedAccidents).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "", rejectedAccidents, w)
	return
}
