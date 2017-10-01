package controllers

import (
	"fmt"
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
	subject = utils.RequestData(r, w)

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

	var categoryID int

	if id, ok := subject["category_id"].(string); ok {
		categoryID, _ = strconv.Atoi(id)
	} else {
		categoryID = int(subject["category_id"].(float64))
	}
	subject["category_id"] = categoryID

	_, err := database.Insert("subjects", subject)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, `{"message": "Subject %s successfully created!"}`, subject["name"])
}

func GetSubjects(w http.ResponseWriter, r *http.Request) {
	db := database.Connect(dbName)

	rows, err := db.Query("SELECT * FROM subjects")
	if err != nil {
		panic(err)
	}

	var allSubjects []Subject

	for {
		if rows.Next() {
			subject := Subject{}
			rows.Scan(&subject.ID, &subject.Name, &subject.CategoryID)
			allSubjects = append(allSubjects, subject)
		} else {
			break
		}
	}

	defer rows.Close()

	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allSubjects)
}

type Subject struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CategoryID int    `json:"category_id"`
}
