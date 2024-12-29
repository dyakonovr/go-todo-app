package api

import (
	"encoding/json"
	"net/http"
	"todo-app/pkg/logger"
)

func NewErrorResponse(w http.ResponseWriter, statusCode int, errorData ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if errorData.Message == "" {
		return
	}
	
	err := json.NewEncoder(w).Encode(errorData)

	if err != nil {
		logger.Errorf("Error while creating error response: %v", err)
		return
	}
}
