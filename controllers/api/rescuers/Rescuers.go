package rescuers

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

func ListOfRescuersApi(w http.ResponseWriter, r *http.Request) {
	var rescuers []structs.Rescuers
	err := models.DB.Raw(`
		SELECT
		    id, 
			name,
			address,
			contact
		FROM
			apps_rescueteamcontacts
	`).Scan(&rescuers).Error

	if err != nil {
		sentry.CaptureException(err)
		// Return an error response if query fails
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Check if user with provided email exists
	if len(rescuers) == 0 {
		message := "Rescuers contact information not found"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}

	// Send success response
	utils.SendSuccessResponse(http.StatusOK, "Rescuer's contact information", rescuers, w)
}

func SelectRescuerApi(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Read and parse request body into SignInInput struct
	var payload structs.SelectRescuer
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
		INSERT INTO apps_accident_rescuer(
			activity_history_id, 
			rescueteamcontacts_id, 
			responders_name,
			notes
		)
		VALUES (?, ?, ?, ?)
	`, payload.ActivityHistoryID, payload.RescuerID, payload.RespondersName, payload.Notes).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	updateActHistory := models.DB.Exec(`
		UPDATE apps_activityhistory
		SET status_report = ? 
		WHERE id = ?
	`, "in-progress", payload.ActivityHistoryID).Error

	if updateActHistory != nil {
		sentry.CaptureException(updateActHistory)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, updateActHistory.Error(), w)
		return
	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully selected rescuer", response, w)
	return
}

func GetRescuerInformation(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	accidentRescuer := queryValues.Get("accident-rescuer")

	var details []structs.RescuerInformationDetails

	err := models.DB.Raw(
		`
			SELECT r.name, r.address, r.contact, a.responders_name, a.notes
			FROM apps_rescueteamcontacts as r
			JOIN apps_accident_rescuer as a ON r.id = a.rescueteamcontacts_id
			WHERE a.activity_history_id = ?
			`, accidentRescuer).Scan(&details).Error

	if len(details) == 0 {
		var empty []interface{}
		utils.SendSuccessResponse(http.StatusOK, "No available assigned rescuers", empty, w)
		return
	}

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved rescuer's details", details, w)
	return
}

func AddNewRescuerApi(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Read and parse request body into SignInInput struct
	var payload structs.AddRescuers
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
		INSERT INTO apps_rescueteamcontacts(
			name, 
			address, 
			contact
		)
		VALUES (?, ?, ?)
	`, payload.Name, payload.Address, payload.Contact).Error

	if execQuery != nil {
		sentry.CaptureException(execQuery)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, execQuery.Error(), w)
		return
	}

	var empty []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "successfully added new rescuers", empty, w)
	return
}
