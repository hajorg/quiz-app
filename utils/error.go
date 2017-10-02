package utils

import (
	"encoding/json"
	"net/http"
)

// Error struct to return a json error
type Error struct {
	Error string `json:"error"`
}

// BadRequest sends a json error to the client with 400 status
func BadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	error := Error{
		Error: err.Error(),
	}
	json.NewEncoder(w).Encode(error)
}

// NotFound sends a json error to the client with a 404 status
func NotFound(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	error := Error{
		Error: err.Error(),
	}
	json.NewEncoder(w).Encode(error)
}

// UnauthorizedError sends a json error to the client with 401 status
func UnauthorizedError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	error := Error{
		Error: message,
	}
	json.NewEncoder(w).Encode(error)
	return
}
