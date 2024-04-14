package utils

import (
	"encoding/json"
	"iotsafedriveapi/structs"
	"net/http"
	"reflect"
)

// SendSuccessResponse sends a success response with the provided data and status code.
func SendSuccessResponse(statusCode int, message string, data interface{}, w http.ResponseWriter) {
	var responseData interface{} // Declare responseData variable outside of if-else block

	checkType := reflect.TypeOf(data).Kind() == reflect.Slice // Check if data is a slice

	if !checkType {
		// Create a slice containing the provided data
		responseData = []interface{}{data} // Assign responseData as a slice with data as its single element
	} else {
		responseData = data // Assign responseData as data itself
	}

	// Create a success response with the provided data
	resp := structs.Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       responseData,
	}

	sendJSONResponse(statusCode, resp, w)
}

// SendErrorResponse sends an error response with the provided message and status code.
func SendErrorResponse(statusCode int, message string, w http.ResponseWriter) {
	// Return an error response with no data
	resp := structs.Response{
		StatusCode: statusCode,
		Message:    message,
	}

	sendJSONResponse(statusCode, resp, w)
}

// sendJSONResponse sends a JSON response with the provided status code and data.
func sendJSONResponse(statusCode int, resp structs.Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode the response as JSON and send it
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
