package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
	"strings"
)

// var dbName string

// func init() {
// 	if flag.Lookup("test.v") == nil {
// 		dbName = "quiz"
// 	} else {
// 		dbName = "quiztest"
// 	}
// }

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

// GetQuestion gets a specific question
func GetQuestion(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "/")
	lastPath := urlPath[len(urlPath)-1]

	data := database.GetWhere("questions", []map[string]interface{}{
		{
			"id": lastPath,
		},
	})

	if len(data) == 0 {
		utils.NotFound(w, errors.New("Question "+lastPath+" does not exist"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(data[0])
}

// GetQuestions gets questions including there options
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	db := database.Connect(dbName)
	defer db.Close()

	stmt, err := db.Query(`SELECT
		questions.*,
		answers.id AS answer_id,
		answers.question_id,
		answers.content AS answer_content,
		answers.correct
		FROM questions
		INNER JOIN answers
		ON questions.id = answers.question_id
	`)
	if err != nil {
		fmt.Println(err)
		utils.NotFound(w, errors.New("Error occurred"))
		return
	}

	// get table columns
	// columns, err := stmt.Columns()
	// if err != nil {
	// 	panic(err)
	// }

	var q Question
	var opt Options
	allQuestions := make(map[int]Question)

	for {
		if stmt.Next() {
			err = stmt.Scan(&q.ID, &q.SubjectID, &q.Content, &opt.AnswerID, &opt.QuestionID, &opt.AnswerContent, &opt.Correct)
			if err != nil {
				panic(err)
			}

			if _, test := allQuestions[q.ID]; test == true {
				tempAns := allQuestions[q.ID].Answers
				tempAns = append(tempAns, opt)
				allQuestions[q.ID] = Question{
					ID:        q.ID,
					SubjectID: q.SubjectID,
					Content:   q.Content,
					Answers:   tempAns,
				}
			} else {
				allQuestions[q.ID] = Question{
					ID:        q.ID,
					SubjectID: q.SubjectID,
					Content:   q.Content,
					Answers: []Options{
						opt,
					},
				}
			}
		} else {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	data := make([]Question, 0)
	for _, value := range allQuestions {
		data = append(data, value)
	}
	json.NewEncoder(w).Encode(data)

}

type Question struct {
	ID        int       `json:"id"`
	SubjectID int       `json:"subjectID"`
	Content   string    `json:"content"`
	Answers   []Options `json:"answers"`
}

type Options struct {
	AnswerID      int    `json:"answerID"`
	QuestionID    int    `json:"questionID"`
	AnswerContent string `json:"answerContent"`
	Correct       bool   `json:"correct"`
}
