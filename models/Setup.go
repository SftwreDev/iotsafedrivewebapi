package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	fmt.Println("Connecting to database")

	// Load the CA certificate
	caCertPath := "ca-certificate.crt"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=verify-ca sslrootcert=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		caCertPath,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
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
