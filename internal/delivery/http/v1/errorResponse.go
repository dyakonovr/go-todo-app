package v1

import (
	"encoding/json"
	"net/http"
	"todo-app/pkg/logger"
)

type dataResponse struct {
	Data        interface{} `json:"data"`
	Count       int         `json:"count"`
	CurrentPage int         `json:"currentPage"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, errorData errorResponse) {
	logger.Error(errorData.Message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	err := json.NewEncoder(w).Encode(errorData)

	if err != nil {
		logger.Errorf("Error while creating error response: %v", err)
		return
	}
}

