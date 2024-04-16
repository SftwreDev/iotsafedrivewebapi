package settings

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func CorsSettings(router *mux.Router) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:8080",
			"http://127.0.0.1:8081",
			"http://192.168.1.19:8081",
			"https://*.ngrok-free.app",
			"https://*.netlify.app",
			"https://iotsafedrive.com",
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},         // Add "Authorization" to the list of allowed headers
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "PATCH"}, // Adjust allowed methods as needed
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	// Insert the middleware
	handler := c.Handler(router)
	return handler
}
