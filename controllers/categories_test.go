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

func TestCreateCategorySuccess(t *testing.T) {
	db := utils.DbTestInit()
	reader := strings.NewReader(`{"title": "testTitle", "description": "testing things"}`)
	req, err := http.NewRequest("POST", "category", reader)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.CreateCategory)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusCreated, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
