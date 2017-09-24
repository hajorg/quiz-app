package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
	"strconv"
)

// CreateSubject creates a new subject
func CreateSubject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// holds form data
	subject := map[string]interface{}{
		"name":        "",
		"category_id": "",
	}

	if len(r.Form) > 0 {
		subject["category_id"] = r.FormValue("category_id")
		subject["name"] = r.FormValue("name")
	} else {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Fprintln(w, err)
		}
		defer r.Body.Close()
		json.Unmarshal(body, &subject)
	}

	validationError := validation.Validator(w, subject, map[string](map[string]string){
		"name": {
			"required": "1",
			"min":      "2",
			"max":      "50",
		},
		"category_id": {
			"required": "1",
		},
	})
	if validationError == false {
		return
	}

	db := database.Connect(dbName)
	defer db.Close()

	smt, err := db.Prepare("INSERT INTO subjects(id, name, category_id) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer smt.Close()
	categoryID, _ := strconv.Atoi(subject["category_id"].(string))
	_, err = smt.Exec(nil, subject["name"], categoryID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, `{"message": "Subject %s successfully created!"}`, subject["name"])
}
