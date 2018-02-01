package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
)

// Results computes result for a particular result
func Results(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := utils.RequestData(r, w)
	db := database.Connect(dbName)
	defer db.Close()

	var currentAnswer Options
	// results = make([]map[])

	for key, val := range data {
		for item := range val {
			stmt, err := db.Prepare("SELECT * FROM answers WHERE question_id = ? AND id = ?")
			defer stmt.Close()
			if err != nil {
				panic(err)
			}

			row, err := stmt.Query(key, val[item])
			defer row.Close()
			if err != nil {
				panic(err)
			}
			if row.Next() {
				row.Scan(
					&currentAnswer.AnswerID,
					&currentAnswer.QuestionID,
					&currentAnswer.AnswerContent,
					&currentAnswer.Correct,
				)
				if currentAnswer.Correct {
					fmt.Println("yes")
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// type Options struct {
// 	AnswerID      int    `json:"answerID"`
// 	QuestionID    int    `json:"questionID"`
// 	AnswerContent string `json:"answerContent"`
// 	Correct       bool   `json:"correct"`
// }
