package routes_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/quiz-app/routes"
	"github.com/quiz-app/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	db      *sql.DB
	err     error
	baseURL = "/api/v1"
)

func TestIndexRoute(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code; got %v but wanted %v", status, http.StatusOK)
	}
}

func TestInvalidRoute(t *testing.T) {
	db := utils.DbTestInit()
	req, err := http.NewRequest("GET", baseURL+"/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code; got %v but wanted %v", status, http.StatusNotFound)
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestUserRouteBadEmailFail(t *testing.T) {
	db := utils.DbTestInit()
	newJson := `{"username": "jameskd", "email": "jamesjd@gmail.com", "password": ""}`
	reader := strings.NewReader(newJson)
	reg, err := http.NewRequest("POST", baseURL+"/user", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(rr, reg)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusUnprocessableEntity, status)
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestUserRouteSuccess(t *testing.T) {
	db := utils.DbTestInit()
	newJson := `{"username": "james", "email": "james@gmail.com", "password": "password"}`
	reader := strings.NewReader(newJson)
	reg, err := http.NewRequest("POST", baseURL+"/user", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(rr, reg)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusCreated, status)
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestCreateCategoryNoTokenFail(t *testing.T) {
	db := utils.DbTestInit()
	newJson := `{"title": "general", "description": "General stuff"}`
	reader := strings.NewReader(newJson)
	req, err := http.NewRequest("POST", baseURL+"/category", reader)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusUnauthorized, status)
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestCreateCategoryPermissionFail(t *testing.T) {
	godotenv.Load()
	db := utils.DbTestInit()
	newJson := `{"title": "general", "description": "General stuff"}`
	reader := strings.NewReader(newJson)
	req, err := http.NewRequest("POST", baseURL+"/category", reader)
	if err != nil {
		t.Fatal(err)
	}
	token := utils.CreateToken(100, 2, "james")
	req.Header.Set("Authorization", token)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusForbidden {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusForbidden, status)
	}

	var body map[string]string
	json.NewDecoder(res.Body).Decode(&body)
	if body["error"] != "You are not permitted to perform this action" {
		t.Errorf("wrong status code: expected `%s` but got %v", "You are not permitted to perform this action", body["error"])
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestCreateCategoryRouteSuccess(t *testing.T) {
	godotenv.Load()
	db := utils.DbTestInit()
	newJson := `{"title": "general", "description": "General stuff"}`
	reader := strings.NewReader(newJson)
	req, err := http.NewRequest("POST", baseURL+"/category", reader)
	if err != nil {
		t.Fatal(err)
	}
	token := utils.CreateToken(100, 1, "james")
	req.Header.Set("Authorization", token)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusCreated, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestCreateSubjectRouteSuccess(t *testing.T) {
	db := utils.DbTestInit()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO categories(id, title, description) VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(nil, "general", "General stuff")
	if err != nil {
		panic(err)
	}
	reader := strings.NewReader(`{"name": "testTitle", "category_id": "1"}`)
	req, _ := http.NewRequest("POST", baseURL+"/subject", reader)
	res := httptest.NewRecorder()

	token := utils.CreateToken(100, 1, "james")
	req.Header.Set("Authorization", token)

	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(res, req)
	t.Log(res)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusCreated, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestCreateQuestionRouteSuccess(t *testing.T) {
	db := utils.DbTestInit()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO categories(id, title, description) VALUES(?, ?, ?)")
	checkError(err)

	defer stmt.Close()
	result, err := stmt.Exec(nil, "general", "General stuff")
	checkError(err)
	lastID, _ := result.LastInsertId()

	stmt, err = db.Prepare("INSERT INTO subjects(id, name, category_id) VALUES(?, ?, ?)")
	checkError(err)
	_, err = stmt.Exec(nil, "English", lastID)
	checkError(err)

	reader := strings.NewReader(`{"subject_id": ` + fmt.Sprint(lastID) + `, "content": "What is a noun?"}`)
	req, _ := http.NewRequest("POST", baseURL+"/question", reader)
	res := httptest.NewRecorder()

	token := utils.CreateToken(100, 1, "james")
	req.Header.Set("Authorization", token)

	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(res, req)
	t.Log(res)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusCreated, status)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS quiztest")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
