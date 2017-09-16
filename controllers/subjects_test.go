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

func TestCreateSubjectSuccess(t *testing.T) {
	db := utils.DbTestInit()
	reader1 := strings.NewReader(`{"title": "testTitle", "description": "testing things"}`)
	req1, err := http.NewRequest("POST", "category", reader1)
	if err != nil {
		panic(err)
	}

	res1 := httptest.NewRecorder()

	handler1 := http.HandlerFunc(controllers.CreateCategory)
	handler1.ServeHTTP(res1, req1)

	reader := strings.NewReader(`{"name": "testTitle", "category_id": "1"}`)
	req, err := http.NewRequest("POST", "subjects", reader)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.CreateSubject)
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
