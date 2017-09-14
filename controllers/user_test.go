package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"quiz-app/controllers"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUserHandlerCreateFail(t *testing.T) {
	_ = DbTestInit()
	data := `{"username": "johny", "email": "johny@yahoo.com", "password": ""}`
	reader := strings.NewReader(data)
	req, _ := http.NewRequest("POST", "user", reader)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)
	t.Log(rr.Code)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusUnprocessableEntity, status)
	}
}

func TestUserHandlerCreateSuccess(t *testing.T) {
	_ = DbTestInit()
	data := `{"username": "john", "email": "john@yahoo.com", "password": "mypassword"}`
	reader := strings.NewReader(data)
	req, _ := http.NewRequest("POST", "user", reader)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)
	t.Log(rr.Code, rr.Body)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusCreated, status)
	}
}

func TestUserHandlerLoginSuccess(t *testing.T) {
	DbTestInit()
	data := `{"username": "john", "password": "mypassword"}`
	reader := strings.NewReader(data)
	req, _ := http.NewRequest("POST", "login", reader)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	t.Log(rr.Code, rr.Body)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusOK, status)
	}
}

func TestUserHandlerLoginFail(t *testing.T) {
	db := DbTestInit()
	data := `{"username": "johnkl", "password": "mypassword"}`
	reader := strings.NewReader(data)
	req, _ := http.NewRequest("POST", "login", reader)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	t.Log(rr.Code, rr.Body)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusUnauthorized, status)
	}
	_, err := db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	db.Close()
}
