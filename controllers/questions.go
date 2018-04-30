package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/quiz-app/database"
	"github.com/quiz-app/utils"
	"github.com/quiz-app/validation"
)

// CreateQuestion creates a single question
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	question := utils.RequestData(r, w)

	valid := validation.Validator(w, question, map[string](map[string]string){
		"subject_id": {
			"required": "1",
		},
		"content": {
			"required": "1",
		},
	})
	if valid == false {
		return
	}

	_, err := database.Insert("questions", question)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(question)
}
