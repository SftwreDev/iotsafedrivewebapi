package main

import (
	"encoding/json"
	"fmt"
	"io"
	"iotsafedriveapi/models"
	"iotsafedriveapi/utils"
	"os"
	"time"
)

type SuperUser struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Address        string `json:"address"`
	Contact        string `json:"contact"`
	Email          string `json:"email"`
	DeviceID       string `json:"device_id"`
	ProfilePicture string `json:"profile_picture"`
	DateJoined     string `json:"date_joined"`
	Password       string `json:"password"`

	IsActive         string `json:"is_active"`
	IsOnboardingDone string `json:"is_onboarding_done"`
	IsStaff          string `json:"is_staff"`
	IsSuperuser      string `json:"is_superuser"`
}

func createSuperUser() {
	// Open our jsonFile
	fmt.Println("Opening json file")
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	fmt.Println("Reading json file")
	byteValue, _ := io.ReadAll(jsonFile)

	var input []SuperUser
	json.Unmarshal([]byte(byteValue), &input)

	for _, superuser := range input {
		// Hash the password securely
		hashedPassword, err := utils.HashPassword(superuser.Password)
		if err != nil {
			fmt.Printf("Error hashing password : %s", err.Error())
		}

		fmt.Printf("Checking if superuser account exists: %s \n", superuser.Email)
		var actor []models.AppsUser

		result := models.DB.Raw(`
		SELECT 
		    first_name, 
		  	last_name, 
		  	email
		FROM
		    apps_user
		WHERE
		    email = ?
	`, superuser.Email).Scan(&actor).Error

		if len(actor) != 0 {
			fmt.Printf("Superuser account already exists: %s \nExiting now... \n", superuser.Email)
		} else {

			fmt.Println("Creating superuser account now")

			result = models.DB.Exec(`
			INSERT INTO 
				apps_user(
						  first_name, 
						  last_name, 
						  email, 
						  password, 
						  is_superuser, 
						  is_staff, 
						  is_onboarding_done, 
						  device_id, 
						  is_active, 
						  date_joined, 
						  profile_picture)
			VALUES 
				(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, superuser.FirstName, superuser.LastName, superuser.Email, hashedPassword, "true", "false",
				"true", superuser.DeviceID, "true", time.Now(), superuser.ProfilePicture).Error

			if result != nil {
				fmt.Printf("Error creating account : %s \n", result)
			} else {
				fmt.Println("Done creating superuser account")
			}
		}
	}

}

//func init() {
//
//	fmt.Println("Server initializing...")
//
//	// Initialize env
//	err := godotenv.Load(".env")
//	if err != nil {
//		return
//	}
//
//	// Initialize models
//	models.ConnectDatabase()
//
//	createSuperUser()
//}
