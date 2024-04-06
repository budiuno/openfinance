package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, errorMessage ErrorResponse, statusCode int) {
	// Set Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code
	w.WriteHeader(statusCode)

	// Encode the error message as JSON
	jsonData, err := json.Marshal(errorMessage)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Write(jsonData)
}

func RespondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	// Set Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code
	w.WriteHeader(statusCode)

	// Encode the data as JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Write(jsonData)
}
