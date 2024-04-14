package auth

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"log"
	"net/http"
)

func ObtainNewToken(w http.ResponseWriter, r *http.Request) {

	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Read and parse request body into RefreshToken struct
	var input structs.RefreshToken
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Validate input using validator package
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		// Return a validation error response
		utils.SendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	refreshClaims := utils.ParseRefreshToken(input.RefreshToken)

	if refreshClaims.Valid() != nil {
		input.RefreshToken, err = utils.NewRefreshToken(*refreshClaims)
		if err != nil {
			log.Fatal("error creating refresh token")
		}
	}

	fmt.Println(input.RefreshToken)

}
