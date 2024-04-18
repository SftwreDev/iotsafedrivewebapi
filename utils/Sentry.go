package utils

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func InitializeSentry() {

	fmt.Println("Initializing Sentry")

	// Initialize sentry config
	sentryError := sentry.Init(sentry.ClientOptions{
		Dsn: "https://e1b7fafa040d80ed75c87983369aeb7e@o4507105149190144.ingest.us.sentry.io/4507105174290432",

		Debug: false,
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
