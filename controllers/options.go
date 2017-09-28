package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
	"strings"
)

// CreateOptions creates a single question
func CreateOption(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	option := map[string]interface{}{}
	if len(r.Form) > 0 {
		for key, val := range r.Form {
			option[key] = strings.Join(val, "")
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		json.Unmarshal(body, &option)
		defer r.Body.Close()
	}

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