package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/quiz-app/controllers"
	"github.com/quiz-app/utils"

	_ "github.com/go-sql-driver/mysql"
)

func TestCreateOptionSuccess(t *testing.T) {
	db := utils.DbTestInit()

	stmt, err := db.Prepare("INSERT INTO categories(id, title, description) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	result, err := stmt.Exec(nil, "general", "General stuff")
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare("INSERT INTO subjects(id, category_id, name) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	lastID, _ := result.LastInsertId()

	result, err = stmt.Exec(nil, lastID, "Maths")
	if err != nil {
		panic(err)
	}

	lastID, _ = result.LastInsertId()

	stmt, err = db.Prepare("INSERT INTO questions(id, subject_id, content) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(nil, lastID, "What is 5*5?")
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(`{"question_id": ` + fmt.Sprint(lastID) + `, "content": 25, "correct": true}`)
	req, err := http.NewRequest("POST", "questions", reader)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.CreateOption)
	handler.ServeHTTP(res, req)
	t.Log(res.Body)
	if status := res.Code; status != http.StatusCreated {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusCreated, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
