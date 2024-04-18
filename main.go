package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"iotsafedriveapi/models"
	"iotsafedriveapi/routes"
	"iotsafedriveapi/utils"
)

func main() {

	fmt.Println("Starting server...")
	// Initialize env
	err := godotenv.Load(".env")
	if err != nil {
		return
	}

	// Initialize Sentry config

	utils.InitializeSentry()

	// Initialize models
	models.ConnectDatabase()
	// models.InitialMigration()

	// // Initialize routers
	routes.InitializeRouter()

}
