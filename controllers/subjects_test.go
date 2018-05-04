package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/quiz-app/controllers"
	"github.com/quiz-app/utils"

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
func TestGetSubjectsSuccess(t *testing.T) {
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

	reader := strings.NewReader(`{"name": "testTitle", "category_id": "1"}`)
	req, err := http.NewRequest("GET", "subject", reader)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.GetSubjects)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusOK, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestGetSubjectSuccess(t *testing.T) {
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

	reader := strings.NewReader(`{"name": "testTitle", "category_id": "1"}`)
	req, err := http.NewRequest("GET", "subject/"+strconv.Itoa(int(lastID)), reader)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.GetSubject)
	handler.ServeHTTP(res, req)
	t.Log(res.Body)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusOK, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestGetSubjectFail(t *testing.T) {
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

	reader := strings.NewReader(`{"name": "testTitle", "category_id": "1"}`)
	req, err := http.NewRequest("GET", "subject/100", reader)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(controllers.GetSubject)
	handler.ServeHTTP(res, req)
	t.Log(res.Body)
	if status := res.Code; status != http.StatusNotFound {
		t.Errorf("Error occurred. Expected %v but got %v status code", http.StatusNotFound, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
