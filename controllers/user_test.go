package controllers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"quiz-app/controllers"
	"quiz-app/database"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS quiztest")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE quiztest")
	if err != nil {
		panic(err)
	}
	database.CreateDatabase("quiztest")
}

func TestUserHandlerCreateFail(t *testing.T) {
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

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/quiztest")
	db.Exec("DROP DATABASE IF EXISTS quiztest")
	db.Close()
}
