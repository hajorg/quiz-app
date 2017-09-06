package utils

import (
	"encoding/json"
	"net/http"
)

// Error struct to return a json error
type Error struct {
	Error string `json:"error"`
}

func BadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	error := Error{
		Error: err.Error(),
	}
	json.NewEncoder(w).Encode(error)
}
