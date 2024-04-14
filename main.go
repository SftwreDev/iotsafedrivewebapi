package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"iotsafedriveapi/models"
	"iotsafedriveapi/routes"
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

}
