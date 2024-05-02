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
					role,
					date_joined
				FROM apps_user;
			`,
	).Scan(&users).Error

	if err != nil {
		sentry.CaptureException(err)
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	// Iterate over users and modify role field
	for i, user := range users {
		// Convert role to human-readable format
		users[i].Role = formatRole(user.Role)
		users[i].DeviceID = formatDeviceID(user.Role, user.DeviceID)
	}

	utils.SendSuccessResponse(http.StatusOK, "Successfully retrieved all users", users, w)
	return
}

// function to format role to human-readable format
func formatRole(role string) string {
	switch role {
	case "super_admin":
		return "Super Admin"
	case "user":
		return "User"
	case "rescuer":
		return "Rescuer"

	default:
		return role // return unchanged if no match
	}
}

func formatDeviceID(role string, deviceID string) string {
	switch role {
	case "super_admin":
		return "N/A"
	case "rescuer":
		return "N/A"

	default:
		return deviceID // return unchanged if no match
	}
}
