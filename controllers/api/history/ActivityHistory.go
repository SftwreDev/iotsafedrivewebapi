package history

import (
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
				WHERE h.status_report IN ('pending', 'in-progress')
				AND h.status = 'SMS Sent'
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
