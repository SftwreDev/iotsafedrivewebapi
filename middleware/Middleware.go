package middleware

import (
	"context"
	"iotsafedriveapi/structs"
	"iotsafedriveapi/utils"
	"net/http"
	"strings"
)

// ValidateToken Middleware function to check if the request has a valid token
func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the header is empty or doesn't start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SendErrorResponse(http.StatusUnauthorized, "Unauthorized: Missing or malformed token", w)
			return
		}

		// Extract the token from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate if access_token is valid or expired
		userClaims, err := utils.ParseAccessToken(tokenString)

		// Check if there's an error parsing the token or if it's invalid
		if err != nil || userClaims.StandardClaims.Valid() != nil {
			utils.SendErrorResponse(http.StatusUnauthorized, "Unauthorized: Expired or invalid token", w)
			return
		}

		// Set parsed userClaims
		r = SetUserClaimsData(
			userClaims.ID, userClaims.Email, userClaims.FirstName, userClaims.LastName, r,
		)

		// Token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}

func SetUserClaimsData(id uint, email string, firstName string, lastName string, r *http.Request) *http.Request {
	userClaim := &structs.UserClaims{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	ctx := context.WithValue(r.Context(), "userClaims", userClaim)
	r = r.WithContext(ctx)

	return r
}
