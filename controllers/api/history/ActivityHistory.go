package history

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
)

func GetAllActivityHistoryApi(w http.ResponseWriter, r *http.Request) {

	var history []structs.ActivityHistory

	err := models.DB.Raw(
		`
				SELECT
					h.id, h.timestamps, h.location, h.latitude, h.longitude, h.status, h.status_report,
					CONCAT(u.first_name, ' ', u.last_name) AS owner
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
					CONCAT(u.first_name, ' ', u.last_name) AS owner
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
					CONCAT(u.first_name, ' ', u.last_name) AS owner
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
