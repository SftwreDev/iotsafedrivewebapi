package auth

import (
	"iotsafedriveapi/utils"
	"net/http"
)

func VerifyTokenApi(w http.ResponseWriter, r *http.Request) {
	// Get userClaims data from request context
	_, ok := utils.GetUserClaimsContext(r)
	if !ok {
		message := "User not found"
		// Return a not found error response
		utils.SendErrorResponse(http.StatusNotFound, message, w)
		return
	}
	var response []string
	// Send success response
	utils.SendSuccessResponse(http.StatusOK, "Token verified", response, w)
	return
}
