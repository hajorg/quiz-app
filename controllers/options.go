package controllers

import (
	"encoding/json"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
)

// CreateOptions creates a single question
func CreateOption(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	option := utils.RequestData(r, w)

	valid := validation.Validator(w, option, map[string](map[string]string){
		"question_id": {
			"required": "1",
		},
		"content": {
			"required": "1",
		},
	})
	if valid == false {
		return
	}

	_, err := database.Insert("answers", option)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(option)
}
