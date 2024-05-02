package analytics

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
	"strings"
)

func DataAnalyticsApi(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	var payload structs.AnalyticsFilter
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

	filter := payload.Filter
	value := payload.Value

	var whereClause string
	var activityHistorySqlQuery string
	var accidentRescueSqlQuery string
	var smsVsFalseSqlQuery string
	var accidentRescuerSqlQuery string
	var onGoingVsClosedVsPendingSqlQuery string

	activityHistoryQuery := `
		SELECT
			COUNT(*) AS total_activity,
			SUM(CASE WHEN status = 'SMS Sent' THEN 1 ELSE 0 END) AS total_sms_sent,
			SUM(CASE WHEN status = 'False Alarm' THEN 1 ELSE 0 END) AS total_false_alarm,
			SUM(CASE WHEN status = 'SMS Not Sent' THEN 1 ELSE 0 END) AS total_sms_not_sent
		FROM apps_activityhistory
	`

	accidentRescueQuery := `
		SELECT COUNT(*) as total_responded_accident
		FROM apps_accident_rescuer
	`

	smsVsFalseQuery := `
		SELECT
			DATE_TRUNC('minute', timestamps) AS timestamps,
			SUM(CASE WHEN status = 'SMS Sent' THEN 1 ELSE 0 END) AS total_sms_sent,
			SUM(CASE WHEN status = 'False Alarm' THEN 1 ELSE 0 END) AS total_false_alarm,
			SUM(CASE WHEN status = 'SMS Not Sent' THEN 1 ELSE 0 END) AS total_sms_not_sent
		FROM apps_activityhistory
	`

	accidentRescuerQuery := `
		SELECT
			DATE_TRUNC('minute', timestamps) AS timestamps,
			Count(*) as total_rescued_accident
		FROM apps_accident_rescuer
	`

	onGoingVsClosedVsPendingQuery := `
		SELECT
			CASE WHEN total_count > 0 THEN ROUND((total_pending * 100.0 / NULLIF(total_count, 0)), 2) ELSE 0 END AS pending_percentage,
			CASE WHEN total_count > 0 THEN ROUND((total_closed * 100.0 / NULLIF(total_count, 0)), 2) ELSE 0 END AS closed_percentage,
			CASE WHEN total_count > 0 THEN ROUND((total_on_going * 100.0 / NULLIF(total_count, 0)), 2) ELSE 0 END AS on_going_percentage
		FROM (
			 	SELECT
					 COUNT(*) FILTER (WHERE status_report = 'pending') AS total_pending,
					 COUNT(*) FILTER (WHERE status_report = 'closed') AS total_closed,
					 COUNT(*) FILTER (WHERE status_report = 'on-going') AS total_on_going,
					 COUNT(*) AS total_count
			 	FROM apps_activityhistory		
		`

	subQueryCloseTag := `) AS subquery;`

	groupByPerMinute := `
		GROUP BY DATE_TRUNC('minute', timestamps)
	`

	switch filter {
	case "yearly":
		whereClause = fmt.Sprintf("WHERE EXTRACT(YEAR FROM timestamps) = %s", value)
		activityHistorySqlQuery = activityHistoryQuery + whereClause
		accidentRescueSqlQuery = accidentRescueQuery + whereClause
		smsVsFalseSqlQuery = smsVsFalseQuery + whereClause + groupByPerMinute
		accidentRescuerSqlQuery = accidentRescuerQuery + whereClause + groupByPerMinute
		onGoingVsClosedVsPendingSqlQuery = onGoingVsClosedVsPendingQuery + whereClause + subQueryCloseTag
	case "monthly":
		month := strings.Split(value, "/")
		whereClause = fmt.Sprintf("WHERE EXTRACT(MONTH FROM timestamps) = %s AND EXTRACT(YEAR FROM timestamps) = %s", month[0], month[1])
		activityHistorySqlQuery = activityHistoryQuery + whereClause
		accidentRescueSqlQuery = accidentRescueQuery + whereClause
		smsVsFalseSqlQuery = smsVsFalseQuery + whereClause + groupByPerMinute
		accidentRescuerSqlQuery = accidentRescuerQuery + whereClause + groupByPerMinute
		onGoingVsClosedVsPendingSqlQuery = onGoingVsClosedVsPendingQuery + whereClause + subQueryCloseTag
	case "daily":
		day := strings.Split(value, "/")
		whereClause = fmt.Sprintf("WHERE EXTRACT(MONTH FROM timestamps) = %s AND  EXTRACT(DAY FROM timestamps) = %s AND  EXTRACT(YEAR FROM timestamps) = %s", day[0], day[1], day[2])
		activityHistorySqlQuery = activityHistoryQuery + whereClause
		accidentRescueSqlQuery = accidentRescueQuery + whereClause
		smsVsFalseSqlQuery = smsVsFalseQuery + whereClause + groupByPerMinute
		accidentRescuerSqlQuery = accidentRescuerQuery + whereClause + groupByPerMinute
		onGoingVsClosedVsPendingSqlQuery = onGoingVsClosedVsPendingQuery + whereClause + subQueryCloseTag
	default:
		activityHistorySqlQuery = activityHistoryQuery
		accidentRescueSqlQuery = accidentRescueQuery
		smsVsFalseSqlQuery = smsVsFalseQuery + groupByPerMinute
		accidentRescuerSqlQuery = accidentRescuerQuery + groupByPerMinute
		onGoingVsClosedVsPendingSqlQuery = onGoingVsClosedVsPendingQuery + subQueryCloseTag
	}

	var activityHistoryStats []structs.ActivityHistoryStats
	var accidentRescueStats []structs.AccidentRescueStats
	var smsVsFalseStats []structs.SmsVsFalse
	var accidentRescuer []structs.AccidentRescuer
	var onGoingVsClosedVsPendingStats []structs.OnGoingVsClosedVsPending

	// SQL Query for Activity History
	execQuery := models.DB.Raw(activityHistorySqlQuery).Scan(&activityHistoryStats).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	// SQL Query for Accident Rescue
	execQuery = models.DB.Raw(accidentRescueSqlQuery).Scan(&accidentRescueStats).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	// SQL Query for SMS vs False
	execQuery = models.DB.Raw(smsVsFalseSqlQuery).Scan(&smsVsFalseStats).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	// SQL Query for Accident Rescuer
	execQuery = models.DB.Raw(accidentRescuerSqlQuery).Scan(&accidentRescuer).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	// SQL Query for OnGoing, Closed and Pending accidents
	execQuery = models.DB.Raw(onGoingVsClosedVsPendingSqlQuery).Scan(&onGoingVsClosedVsPendingStats).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response if token generation fails
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	results := make(map[string]interface{})
	results["activity_history"] = activityHistoryStats
	results["accident_rescue"] = accidentRescueStats
	results["sms_vs_false"] = smsVsFalseStats
	results["rescued_accident"] = accidentRescuer
	results["ongoing_vs_closed_pending"] = onGoingVsClosedVsPendingStats

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieve your statistics", results, w)
	return

}
