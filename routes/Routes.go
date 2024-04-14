package routes

import (
	"iotsafedriveapi/controllers/api"
	"iotsafedriveapi/controllers/api/accident_alert"
	"iotsafedriveapi/controllers/api/actor"
	"iotsafedriveapi/controllers/api/auth"
	"iotsafedriveapi/controllers/api/history"
	"iotsafedriveapi/controllers/api/rescuers"
	"iotsafedriveapi/controllers/api/trusted_contacts"
	"iotsafedriveapi/controllers/api/vehicle"
	"iotsafedriveapi/middleware"
	"iotsafedriveapi/settings"
	"log"
	"net/http"
)

func InitializeRouter() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/signup", auth.SignUpApi)
	mux.HandleFunc("POST /api/signin", auth.SignInApi)
	mux.HandleFunc("POST /api/reset-password", auth.ResetPasswordApi)
	mux.HandleFunc("POST /api/token/obtain", auth.ObtainNewToken)
	mux.Handle("GET /api/token/verify", middleware.ValidateToken(http.HandlerFunc(auth.VerifyTokenApi)))
	mux.Handle("POST /api/update-password", middleware.ValidateToken(http.HandlerFunc(auth.UpdatePasswordApi)))

	// Use the ValidateToken middleware for protected routes
	mux.Handle("GET /api/actors", middleware.ValidateToken(http.HandlerFunc(api.ActorGetListApi)))
	mux.Handle("GET /api/actor/profile", middleware.ValidateToken(http.HandlerFunc(actor.GetActorProfileApi)))
	mux.Handle("POST /api/actor/profile/update", middleware.ValidateToken(http.HandlerFunc(actor.UpdateActorProfileApi)))

	mux.Handle("GET /api/rescuers", middleware.ValidateToken(http.HandlerFunc(rescuers.ListOfRescuersApi)))
	mux.Handle("POST /api/rescuers/select", middleware.ValidateToken(http.HandlerFunc(rescuers.SelectRescuerApi)))
	mux.Handle("GET /api/rescuers/details", middleware.ValidateToken(http.HandlerFunc(rescuers.GetRescuerInformation)))
	mux.Handle("POST /api/rescuers/add", middleware.ValidateToken(http.HandlerFunc(rescuers.AddNewRescuerApi)))

	mux.Handle("GET /api/vehicle", middleware.ValidateToken(http.HandlerFunc(vehicle.GetUsersVehicleApi)))
	mux.Handle("POST /api/vehicle/add", middleware.ValidateToken(http.HandlerFunc(vehicle.AddVehicleApi)))
	mux.Handle("POST /api/vehicle/update", middleware.ValidateToken(http.HandlerFunc(vehicle.UpdateVehicleApi)))
	mux.Handle("GET /api/vehicle/all", middleware.ValidateToken(http.HandlerFunc(vehicle.GetAllVehicleApi)))

	mux.Handle("GET /api/trusted-contacts", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.ListTrustedContactsApi)))
	mux.Handle("GET /api/trusted-contacts/all", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.GetAllTrustedContactsApi)))
	mux.Handle("POST /api/trusted-contacts/add", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.AddTrustedContactsApi)))
	mux.Handle("POST /api/trusted-contacts/update", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.UpdateTrustedContactsApi)))

	mux.Handle("GET /api/activity-history/all", middleware.ValidateToken(http.HandlerFunc(history.GetAllActivityHistoryApi)))
	mux.Handle("GET /api/activity-history/pending", middleware.ValidateToken(http.HandlerFunc(history.GetPendingActivityHistoryApi)))
	mux.Handle("DELETE /api/activity-history/close", middleware.ValidateToken(http.HandlerFunc(history.CloseActivityHistoryApi)))
	mux.Handle("GET /api/activity-history", middleware.ValidateToken(http.HandlerFunc(history.GetDetailedActivityHistoryApi)))
	mux.Handle("GET /api/activity-history/latest", middleware.ValidateToken(http.HandlerFunc(history.GetLatestActivityHistoryApi)))

	mux.Handle("GET /api/users/all", middleware.ValidateToken(http.HandlerFunc(actor.GetAllUsersApi)))

	mux.Handle("GET /api/accident-alert", middleware.ValidateToken(http.HandlerFunc(accident_alert.GetAllAccidentAlertApi)))

	mux.Handle("POST /api/iot/alerts", http.HandlerFunc(accident_alert.AccidentDetectedApi))
	mux.Handle("GET /api/iot/alerts/check", http.HandlerFunc(accident_alert.GetLatestAccidentAlertApi))
	mux.Handle("POST /api/send-sms", middleware.ValidateToken(http.HandlerFunc(accident_alert.SendSMSApi)))

	// Inserting CORS middleware settings
	handler := settings.CorsSettings(mux)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
