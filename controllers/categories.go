package controllers

import (
	"encoding/json"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
)

// CreateCategory creates new category
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	category := utils.RequestData(r, w)

	valid := validation.Validator(w, category, map[string](map[string]string){
		"title": {
			"required": "1",
			"min":      "2",
			"max":      "50",
		},
		"description": {
			"max": "255",
		},
	})
	if valid == false {
		return
	}

	_, err := database.Insert("categories", category)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(category)
}

// GetCategories gets all categories from the database
func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := database.GetAll("categories")
	json.NewEncoder(w).Encode(data)
}
