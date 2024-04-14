package settings

import (
	"github.com/rs/cors"
	"net/http"
)

func CorsSettings(mux *http.ServeMux) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:8080",
			"http://127.0.0.1:8081",
			"http://192.168.1.19:8081",
			"https://*.ngrok-free.app",
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},         // Add "Authorization" to the list of allowed headers
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "PATCH"}, // Adjust allowed methods as needed
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// Insert the middleware
	handler := c.Handler(mux)
	return handler
}
