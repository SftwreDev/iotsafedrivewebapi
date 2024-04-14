package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// ConnectDatabase establishes a connection to the PostgreSQL database using the
	// provided environment variables for configuration. It initializes the GORM database
	// instance and performs automatic migration for the specified models.
	//
	// This function expects the following environment variables to be set:
	// - DB_HOST: PostgreSQL host address
	// - DB_USER: PostgreSQL username
	// - DB_PASSWORD: PostgreSQL password
	// - DB_NAME: PostgreSQL database name
	// - DB_PORT: PostgreSQL port number
	//
	// Upon successful connection and migration, the global variable DB will be set to
	// the GORM database instance for further use.
	//
	// Example usage:
	//     ConnectDatabase()
	//
	// Note: Ensure that the required GORM models are imported and defined before calling
	// this function.

	fmt.Println("Connecting to database")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	DB = database
	fmt.Println("Successfully connected to database")

}

//func InitialMigration() {
//
//	// Add models here for migration
//	err := DB.AutoMigrate(&AppsUser{})
//
//	if err != nil {
//		return
//	}
//}
