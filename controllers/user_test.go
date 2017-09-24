package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"quiz-app/controllers"
	"quiz-app/database"
	"quiz-app/utils"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUserHandlerCreateFail(t *testing.T) {
	db := utils.DbTestInit()
	database.CreateDatabase("quiztest")
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
	_, err := db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestUserHandlerCreateSuccess(t *testing.T) {
	db := utils.DbTestInit()
	database.CreateDatabase("quiztest")
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
	_, err := db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestUserHandlerLoginSuccess(t *testing.T) {
	db := utils.DbTestInit()
	database.CreateDatabase("quiztest")
	signupData := `{"username": "john", "email": "johny@yahoo.com", "password": "mypassword"}`
	signupReader := strings.NewReader(signupData)
	signupReq, _ := http.NewRequest("POST", "user", signupReader)

	signupRr := httptest.NewRecorder()

	handler1 := http.HandlerFunc(controllers.CreateUser)
	handler1.ServeHTTP(signupRr, signupReq)

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
	_, err := db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestUserHandlerLoginFail(t *testing.T) {
	db := utils.DbTestInit()
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
	defer db.Close()
}
