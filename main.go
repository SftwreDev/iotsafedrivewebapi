package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"iotsafedriveapi/models"
	"iotsafedriveapi/routes"
	"log"
	"strings"
	"time"
)

func main() {

	fmt.Println("Starting server...")
	// Initialize env
	err := godotenv.Load(".env")
	if err != nil {
		return
	}

	// Initialize models
	models.ConnectDatabase()
	// models.InitialMigration()

	// // Initialize routers
	routes.InitializeRouter()

	sentryError := sentry.Init(sentry.ClientOptions{
		Dsn: "https://e1b7fafa040d80ed75c87983369aeb7e@o4507105149190144.ingest.us.sentry.io/4507105174290432",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		Debug: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Here you can inspect/modify non-transaction events (for example, errors) before they are sent.
			// Returning nil drops the event.
			log.Printf("BeforeSend event [%s]", event.EventID)
			return event
		},
		BeforeSendTransaction: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Here you can inspect/modify transaction events before they are sent.
			// Returning nil drops the event.
			if strings.Contains(event.Message, "test-transaction") {
				// Drop the transaction
				return nil
			}
			event.Message += " [example]"
			log.Printf("BeforeSendTransaction event [%s]", event.EventID)
			return event
		},
		// Enable tracing
		EnableTracing: true,
		// Specify either a TracesSampleRate...
		TracesSampleRate: 1.0,
	})
	if sentryError != nil {
		log.Fatalf("sentry.Init: %s", sentryError)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

}
