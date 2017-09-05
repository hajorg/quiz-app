package routes_test

import (
	"net/http"
	"net/http/httptest"
	"quiz-app/routes"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestRoutes(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.Routers)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != 200 {
		t.Errorf("handler returned wrong status code; got %v but wanted %v", status, http.StatusOK)
	}
}
