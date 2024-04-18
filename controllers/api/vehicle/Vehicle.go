package vehicle

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

func AddVehicleApi(w http.ResponseWriter, r *http.Request) {
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

	var input structs.Vehicle
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Validate input using validator package
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		sentry.CaptureException(err)
		// Return a validation error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Insert data to db
	result := models.DB.Exec(`
		INSERT INTO apps_vehicle(brand, model, year_model, plate_no, owner_id) 
		VALUES (?, ?, ?, ?, ?)
		`, input.Brand, input.Model, input.YearModel, input.PlateNo, userID).Error

	if result != nil {
		sentry.CaptureException(result)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error(), w)
		return
	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully added your vehicle", response, w)
	return
}

func GetUsersVehicleApi(w http.ResponseWriter, r *http.Request) {
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

	var vehicle []structs.Vehicle

	err := models.DB.Raw(
		`
				SELECT * FROM apps_vehicle
				WHERE owner_id = ?
			`, userID,
	).Scan(&vehicle).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved your vehicle info", vehicle, w)
	return
}

func UpdateVehicleApi(w http.ResponseWriter, r *http.Request) {
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

	var input structs.Vehicle
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Validate input using validator package
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		sentry.CaptureException(err)
		// Return a validation error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Insert data to db
	result := models.DB.Exec(`
		UPDATE apps_vehicle
		SET
		    brand = ?,
		    model = ?,
		    year_model = ?,
		    plate_no = ?
		WHERE
		    owner_id = ?
		`, input.Brand, input.Model, input.YearModel, input.PlateNo, userID).Error

	if result != nil {
		sentry.CaptureException(err)
		// Return an error response
		utils.SendErrorResponse(http.StatusBadRequest, result.Error(), w)
		return
	}

	var response []interface{}
	utils.SendSuccessResponse(http.StatusCreated, "Successfully added your vehicle", response, w)
	return

}

func GetAllVehicleApi(w http.ResponseWriter, r *http.Request) {

	var vehicle []structs.AllVehicle

	err := models.DB.Raw(
		`
				SELECT
					v.owner_id, v.brand, v.model, v.year_model, v.plate_no,
					CONCAT(u.first_name, ' ', u.last_name) AS owner
				FROM apps_vehicle as v
				INNER JOIN apps_user as u ON u.id = v.owner_id;
			`,
	).Scan(&vehicle).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved your vehicle info", vehicle, w)
	return
}
