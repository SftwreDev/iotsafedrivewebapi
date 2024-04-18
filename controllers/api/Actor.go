package api

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"iotsafedriveapi/models"
	"iotsafedriveapi/structs"
	"net/http"
)

func ActorGetListApi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Set the Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Declare actors models
	var appsUser []models.AppsUser

	// Execute the database query
	err := models.DB.Select("id, first_name, last_name, address, contact, email, device_id, profile_picture, is_active, is_onboarding_done, is_staff, is_superuser").
		WithContext(ctx).
		Find(&appsUser).Error

	if err != nil {
		sentry.CaptureException(err)
		resp := structs.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}
		// Encode the response as JSON and send it
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	// Convert models to response structs
	var response []structs.AppsUser
	for _, user := range appsUser {
		response = append(response, structs.AppsUser{
			ID:               user.ID,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Address:          user.Address,
			Contact:          user.Contact,
			Email:            user.Email,
			DeviceID:         user.DeviceID,
			ProfilePicture:   user.ProfilePicture,
			IsActive:         user.IsActive,
			IsOnboardingDone: user.IsOnboardingDone,
			IsStaff:          user.IsStaff,
			IsSuperuser:      user.IsSuperuser,
		})
	}

	// Create a success response
	resp := structs.Response{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       response,
	}
	// Encode the response as JSON and send it
	_ = json.NewEncoder(w).Encode(resp)
}
