package routes

import (
	"github.com/gorilla/mux"
	"iotsafedriveapi/controllers/api"
	"iotsafedriveapi/controllers/api/accident_alert"
	"iotsafedriveapi/controllers/api/actor"
	"iotsafedriveapi/controllers/api/analytics"
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
	router := mux.NewRouter()

	router.HandleFunc("/api/signup", auth.SignUpApi).Methods("POST")
	router.HandleFunc("/api/signin", auth.SignInApi).Methods("POST")
	router.HandleFunc("/api/reset-password", auth.ResetPasswordApi).Methods("POST")
	router.HandleFunc("/api/token/obtain", auth.ObtainNewToken).Methods("POST")
	router.Handle("/api/token/verify", middleware.ValidateToken(http.HandlerFunc(auth.VerifyTokenApi))).Methods("GET")
	router.Handle("/api/update-password", middleware.ValidateToken(http.HandlerFunc(auth.UpdatePasswordApi))).Methods("POST")
	router.Handle("/api/update-temporary-password", middleware.ValidateToken(http.HandlerFunc(auth.UpdateTemporaryPassword))).Methods("POST")

	// Use the ValidateToken middleware for protected routes
	router.Handle("/api/actors", middleware.ValidateToken(http.HandlerFunc(api.ActorGetListApi))).Methods("GET")
	router.Handle("/api/actor/profile", middleware.ValidateToken(http.HandlerFunc(actor.GetActorProfileApi))).Methods("GET")
	router.Handle("/api/actor/profile/update", middleware.ValidateToken(http.HandlerFunc(actor.UpdateActorProfileApi))).Methods("POST")
	router.Handle("/api/actor/is-password-changed", middleware.ValidateToken(http.HandlerFunc(actor.IsPasswordChangedApi))).Methods("GET")

	router.Handle("/api/rescuers", middleware.ValidateToken(http.HandlerFunc(rescuers.ListOfRescuersApi))).Methods("GET")
	router.Handle("/api/rescuers/select", middleware.ValidateToken(http.HandlerFunc(rescuers.SelectRescuerApi))).Methods("POST")
	router.Handle("/api/rescuers/details", middleware.ValidateToken(http.HandlerFunc(rescuers.GetRescuerInformation))).Methods("GET")
	router.Handle("/api/rescuers/add", middleware.ValidateToken(http.HandlerFunc(rescuers.AddNewRescuerApi))).Methods("POST")

	router.Handle("/api/vehicle", middleware.ValidateToken(http.HandlerFunc(vehicle.GetUsersVehicleApi))).Methods("GET")
	router.Handle("/api/vehicle/add", middleware.ValidateToken(http.HandlerFunc(vehicle.AddVehicleApi))).Methods("POST")
	router.Handle("/api/vehicle/update", middleware.ValidateToken(http.HandlerFunc(vehicle.UpdateVehicleApi))).Methods("POST")
	router.Handle("/api/vehicle/all", middleware.ValidateToken(http.HandlerFunc(vehicle.GetAllVehicleApi))).Methods("GET")

	router.Handle("/api/trusted-contacts", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.ListTrustedContactsApi))).Methods("GET")
	router.Handle("/api/trusted-contacts/all", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.GetAllTrustedContactsApi))).Methods("GET")
	router.Handle("/api/trusted-contacts/add", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.AddTrustedContactsApi))).Methods("POST")
	router.Handle("/api/trusted-contacts/update", middleware.ValidateToken(http.HandlerFunc(trusted_contacts.UpdateTrustedContactsApi))).Methods("POST")

	router.Handle("/api/activity-history/all", middleware.ValidateToken(http.HandlerFunc(history.GetAllActivityHistoryApi))).Methods("GET")
	router.Handle("/api/activity-history/pending", middleware.ValidateToken(http.HandlerFunc(history.GetPendingActivityHistoryApi))).Methods("GET")
	router.Handle("/api/activity-history/close", middleware.ValidateToken(http.HandlerFunc(history.CloseActivityHistoryApi))).Methods("DELETE")
	router.Handle("/api/activity-history", middleware.ValidateToken(http.HandlerFunc(history.GetDetailedActivityHistoryApi))).Methods("GET")
	router.Handle("/api/activity-history/latest", middleware.ValidateToken(http.HandlerFunc(history.GetLatestActivityHistoryApi))).Methods("GET")
	router.Handle("/api/activity-history/forwarded", middleware.ValidateToken(http.HandlerFunc(history.GetForwardedAccidentsApi))).Methods("GET")

	router.Handle("/api/users/all", middleware.ValidateToken(http.HandlerFunc(actor.GetAllUsersApi))).Methods("GET")

	router.Handle("/api/accident-alert", middleware.ValidateToken(http.HandlerFunc(accident_alert.GetAllAccidentAlertApi))).Methods("GET")
	router.Handle("/api/accident-alert/forward", middleware.ValidateToken(http.HandlerFunc(accident_alert.ForwardAccidentApi))).Methods("POST")
	router.Handle("/api/accident-alert/status", middleware.ValidateToken(http.HandlerFunc(accident_alert.CheckIfAccidentIsForwarded))).Methods("GET")
	router.Handle("/api/accident-alert/action", middleware.ValidateToken(http.HandlerFunc(accident_alert.ForwardedAccidentsActions))).Methods("POST")
	router.Handle("/api/accident-alert/accepted", middleware.ValidateToken(http.HandlerFunc(accident_alert.AcceptedAccidents))).Methods("GET")
	router.Handle("/api/accident-alert/rejected", middleware.ValidateToken(http.HandlerFunc(accident_alert.RejectedAccidents))).Methods("GET")

	router.Handle("/api/iot/alerts", http.HandlerFunc(accident_alert.AccidentDetectedApi)).Methods("POST")
	router.Handle("/api/iot/alerts/check", middleware.ValidateToken(http.HandlerFunc(accident_alert.GetLatestAccidentAlertApi))).Methods("GET")
	router.Handle("/api/send-sms", middleware.ValidateToken(http.HandlerFunc(accident_alert.SendSMSApi))).Methods("POST")

	router.Handle("/api/add-account", http.HandlerFunc(actor.AddAccountApi)).Methods("POST")

	router.Handle("/api/analytics/statistics", middleware.ValidateToken(http.HandlerFunc(analytics.DataAnalyticsApi))).Methods("POST")

	// Inserting CORS middleware settings
	handler := settings.CorsSettings(router)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}

}
