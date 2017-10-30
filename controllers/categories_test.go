package controllers_test

import (
	"encoding/json"
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

func TestGetCategorySuccess(t *testing.T) {
	db := utils.DbTestInit()
	_, err := db.Exec("INSERT INTO categories(id, title, description) VALUES(1, 'general', 'General stuff')")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "category", nil)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.GetCategories)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusOK, status)
	}
	var resBody interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)
	result, ok := resBody.([]interface{})
	if !ok {
		panic("Error occurred")
	}

	if length := len(result); length != 1 {
		t.Errorf("Error occurred. Expected length of %v but got %v ", 1, length)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
