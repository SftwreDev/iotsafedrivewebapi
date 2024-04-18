package actor

import (
	"github.com/getsentry/sentry-go"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
)

func GetAllUsersApi(w http.ResponseWriter, r *http.Request) {

	var users []structs.AllUsers

	err := models.DB.Raw(
		`
				SELECT
					first_name,
					last_name,
					email,
					address,
					contact,
					device_id,
					date_joined
				FROM apps_user;
			`,
	).Scan(&users).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved all users", users, w)
	return
}
