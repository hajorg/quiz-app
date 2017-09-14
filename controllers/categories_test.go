package controllers_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"quiz-app/controllers"
	"quiz-app/database"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCreateCategorySuccess(t *testing.T) {
	db := DbTestInit()
	fmt.Println("start2")
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

func DbTestInit() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")
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
	return db
}
