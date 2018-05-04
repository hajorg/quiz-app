package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quiz-app/controllers"
)

func TestIndex(t *testing.T) {
	reg, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Index)

	handler.ServeHTTP(rr, reg)
	if status := rr.Code; status != 200 {
		t.Errorf("handler returned wrong status code; got %v but wanted %v", status, http.StatusOK)
	}
}
