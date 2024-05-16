package utils

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"iotsafedriveapi/models"
	"time"
)

func PerformPendingAccidents(activityHistoryID string) {
	fmt.Println("Background task started")
	time.Sleep(20 * time.Second)

	var status string

	execQuery := models.DB.Raw(`
					SELECT status_report 
					FROM apps_activityhistory 
					WHERE id = ?`, activityHistoryID).Scan(&status).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
	}

	if status == "forwarded" {
		execQuery = models.DB.Exec(`
				UPDATE apps_activityhistory 
				SET status_report = 'pending' 
				WHERE id = ?`, activityHistoryID).Error

		if execQuery != nil {
			sentry.CaptureException(execQuery)
		}

		execQuery = models.DB.Exec(`
			DELETE FROM apps_forwarded_accidents 
			WHERE activity_history_id = ?`, activityHistoryID).Error
		if execQuery != nil {
			sentry.CaptureException(execQuery)
		}
		fmt.Println("Reverting forwarded accidents to base command")
	}

	fmt.Println("Background task completed")
}
