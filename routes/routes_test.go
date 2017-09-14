package routes_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quiz-app/database"
	"quiz-app/routes"
	"quiz-app/utils"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code; got %v but wanted %v", status, http.StatusNotFound)
	}
}

func TestUserRouteBadEmailFail(t *testing.T) {
	newJson := `{"username": "jameskd", "email": "jamesjd@gmail.com", "password": ""}`
	reader := strings.NewReader(newJson)
	reg, err := http.NewRequest("POST", "/user", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(rr, reg)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusUnprocessableEntity, status)
	}
}

func TestUserRouteSuccess(t *testing.T) {
	newJson := `{"username": "james", "email": "james@gmail.com", "password": "password"}`
	reader := strings.NewReader(newJson)
	reg, err := http.NewRequest("POST", "/user", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(rr, reg)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusCreated, status)
	}
}

func TestCreateCategoryNoTokenFail(t *testing.T) {
	newJson := `{"title": "general", "description": "General stuff"}`
	reader := strings.NewReader(newJson)
	req, err := http.NewRequest("POST", "/category", reader)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusUnauthorized, status)
	}
}

func TestCreateCategoryPermissionFail(t *testing.T) {
	godotenv.Load()
	newJson := `{"title": "general", "description": "General stuff"}`
	reader := strings.NewReader(newJson)
	req, err := http.NewRequest("POST", "/category", reader)
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
}

func TestCreateCategorySuccess(t *testing.T) {
	godotenv.Load()
	newJson := `{"title": "general", "description": "General stuff"}`
	reader := strings.NewReader(newJson)
	req, err := http.NewRequest("POST", "/category", reader)
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

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/quiztest")
	db.Exec("DROP DATABASE IF EXISTS quiztest")
	db.Close()
}
