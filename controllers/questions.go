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

// CreateQuestion creates a single question
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	question := map[string]interface{}{}
	if len(r.Form) > 0 {
		for key, val := range r.Form {
			question[key] = strings.Join(val, "")
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		json.Unmarshal(body, &question)
		defer r.Body.Close()
	}

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

	db := database.Connect(dbName)
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO questions(id, subject_id, content) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(nil, question["subject_id"], question["content"])
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(question)
}
