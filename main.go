package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"iotsafedriveapi/models"
	"iotsafedriveapi/routes"
	"log"
)

func main() {

	log.Println("Starting webserver...")
	// Initialize env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Initialize Sentry config

	//utils.InitializeSentry()

	// Initialize models
	models.ConnectDatabase()
	// models.InitialMigration()

	// // Initialize routers
	routes.InitializeRouter()

}
