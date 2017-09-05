package controllers

import (
	"fmt"
	"net/http"
)

// Index The index page
func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to Quizzy!")
}
