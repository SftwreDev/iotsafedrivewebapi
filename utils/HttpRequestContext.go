package utils

import (
	"iotsafedriveapi/structs"
	"net/http"
)

func GetUserClaimsContext(r *http.Request) (*structs.UserClaims, bool) {
	userClaims, ok := r.Context().Value("userClaims").(*structs.UserClaims)
	if !ok {
		return nil, false
	}
	return userClaims, ok
}
