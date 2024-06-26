package history

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

func GetAllActivityHistoryApi(w http.ResponseWriter, r *http.Request) {

	var history []structs.ActivityHistory

	err := models.DB.Raw(
		`
				SELECT
					h.id, h.timestamps, h.location, h.latitude, h.longitude, h.status, h.status_report,
					CONCAT(u.first_name, ' ', u.last_name) AS owner, u.device_id
				FROM apps_activityhistory as h
				INNER JOIN apps_user as u ON u.id = h.user_id;
			`,
	).Scan(&history).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved your vehicle info", history, w)
	return
}

func GetPendingActivityHistoryApi(w http.ResponseWriter, r *http.Request) {

	var history []structs.ActivityHistory

	err := models.DB.Raw(
		`
				SELECT
					h.id, h.timestamps, h.location, h.latitude, h.longitude, h.status, h.status_report,
					CONCAT(u.first_name, ' ', u.last_name) AS owner
				FROM apps_activityhistory as h
				INNER JOIN apps_user as u ON u.id = h.user_id
				WHERE h.status = 'SMS Sent'
				ORDER BY h.timestamps DESC
			`,
	).Scan(&history).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved your vehicle info", history, w)
	return
}

func GetDetailedActivityHistoryApi(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from the request URL
	queryValues := r.URL.Query()

	activityId := queryValues.Get("id")

	var history []structs.ActivityHistory

	err := models.DB.Raw(
		`
				SELECT
					h.id, h.timestamps, h.location, h.latitude, h.longitude, h.status, h.status_report,
					CONCAT(u.first_name, ' ', u.last_name) AS owner, u.device_id
				FROM apps_activityhistory as h
				INNER JOIN apps_user as u ON u.id = h.user_id
				WHERE h.id = ?
				ORDER BY h.timestamps DESC
			`, activityId).Scan(&history).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved detailed activity history", history, w)
	return
}

func CloseActivityHistoryApi(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from the request URL
	queryValues := r.URL.Query()

	activityId := queryValues.Get("id")

	err := models.DB.Exec(
		`
				UPDATE apps_activityhistory
				Set status_report = 'closed'
				WHERE id = ?
			`, activityId).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusOK, "Successfully closed a activity history", response, w)
	return
}

func GetLatestActivityHistoryApi(w http.ResponseWriter, r *http.Request) {

	var history []structs.ActivityHistory

	err := models.DB.Raw(
		`
				SELECT
					h.id, h.timestamps, h.location, h.latitude, h.longitude, h.status, h.status_report,
					CONCAT(u.first_name, ' ', u.last_name) AS owner, u.device_id
				FROM apps_activityhistory as h
						 INNER JOIN apps_user as u ON u.id = h.user_id
				WHERE status_report = 'pending' AND status = 'SMS Sent'
				ORDER BY h.timestamps DESC
				LIMIT 1
			`).Scan(&history).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	if len(history) == 0 {
		message := "No pending accident as of the moment"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusOK, message, w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved latest activity history", history, w)
	return
}

func GetForwardedAccidentsApi(w http.ResponseWriter, r *http.Request) {

	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters from the request URL
	queryValues := r.URL.Query()

	queryType := queryValues.Get("type")

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

	var accidents []structs.ForwardedAccidents
	var sqlQuery string

	query := fmt.Sprintf(`
		SELECT
			fa.id,
			fa.notes,
			fa.status,
			fa.forwarded_by,
			fa.activity_history_id,
			fa.timestamps as forwarded_on,
			ah.location as location,
			ah.timestamps as accident_occurred_on,
			CONCAT(u.first_name, ' ', u.last_name) as victim
		FROM apps_forwarded_accidents as fa
		INNER JOIN apps_activityhistory as ah ON fa.activity_history_id = ah.id
		INNER JOIN apps_user as u ON ah.user_id = u.id
		WHERE fa.rescuer_id = %d
		ORDER BY fa.timestamps DESC
	`, userID)

	if queryType == "limit" {
		sqlQuery = query + " LIMIT 1"
	} else {
		sqlQuery = query
	}

	err := models.DB.Raw(sqlQuery).Scan(&accidents).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	if len(accidents) == 0 {
		utils.SendErrorResponse(http.StatusBadRequest, "Forwarded accidents not found", w)
		return
	} else {
		utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved your vehicle info", accidents, w)
		return
	}

}

func GetSpecificActivityHistoryForNarrativeReportApi(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from the request URL
	queryValues := r.URL.Query()

	activityId := queryValues.Get("id")

	var details []structs.AccidentDetails

	err := models.DB.Raw(
		`
				SELECT
					h.id,
					h.timestamps,
					h.location,
					h.latitude,
					h.longitude,
					h.status,
					h.status_report,
					CONCAT(u.first_name, ' ', u.last_name) AS owner,
					u.device_id,
					CONCAT(au.first_name, ' ', au.last_name) AS rescuer
				FROM apps_activityhistory as h
				JOIN apps_user as u ON u.id = h.user_id
				JOIN apps_accepted_accidents aaa on h.id = aaa.activity_history_id
				JOIN apps_user as au ON au.id = aaa.rescuer_id
				WHERE h.id = ?
				ORDER BY h.timestamps DESC
			`, activityId).Scan(&details).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved detailed activity history", details, w)
	return
}

func CreateNarrativeReportApi(w http.ResponseWriter, r *http.Request) {

	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get userClaims data from request context
	userClaims, ok := utils.GetUserClaimsContext(r)
	if !ok {
		message := "User not found"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	userID := userClaims.ID

	// Read and parse request body into SignInInput struct
	var payload structs.CreateNarrativeReport
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

	timestamps := time.Now()

	execQuery := models.DB.Exec(`
		INSERT INTO 
		    apps_narrative_report(activity_history_id, description, reported_by, timestamps)
		VALUES(?, ?, ?, ?)
	`, payload.ActivityHistoryID, payload.Description, strconv.Itoa(int(userID)), timestamps).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusOK, "Successfully forwarded accident", response, w)
	return
}
