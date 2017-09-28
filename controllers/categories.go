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

// CreateCategory creates new category
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	category := map[string]interface{}{}
	if len(r.Form) > 0 {
		for key, val := range r.Form {
			category[key] = strings.Join(val, "")
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		json.Unmarshal(body, &category)
		defer r.Body.Close()
	}

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
