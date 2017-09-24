package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"quiz-app/controllers"
	"quiz-app/utils"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCreateQuestionSuccess(t *testing.T) {
	db := utils.DbTestInit()

	stmt, err := db.Prepare("INSERT INTO categories(id, title, description) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(nil, "general", "General stuff")
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare("INSERT INTO subjects(id, category_id, name) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(nil, 1, "Maths")
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(`{"subject_id": "1", "content": "What is IoT?"}`)
	req, err := http.NewRequest("POST", "subjects", reader)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.CreateQuestion)
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
