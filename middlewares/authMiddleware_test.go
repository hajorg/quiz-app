package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"quiz-app/utils"
	"testing"
)

func TestAuthMiddlewareSuccess(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		panic(err)
	}

	token := utils.CreateToken(100, 1, "james")
	req.Header.Set("Authorization", token)
	res := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	handler.ServeHTTP(res, req)
	t.Log(res.Code)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusOK, status)
	}
}

func TestAuthMiddlewareFail(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		panic(err)
	}

	res := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	handler.ServeHTTP(res, req)
	t.Log(res.Code)
	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusUnauthorized, status)
	}
}

func TestAuthMiddlewareEmptyTokenFail(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJqYW5lIiwicm9sZUlkIjoyLCJleHAiOjE1MDUzOTgzNzF9.Iqz9HBdbuoTm1hB4BX3K7wdZEfwei_n3RinwcAdvmpg")

	res := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	handler.ServeHTTP(res, req)
	t.Log(res.Body)
	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusUnauthorized, status)
	}
}

func TestAuthAdminMiddlewareSuccess(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		panic(err)
	}

	token := utils.CreateToken(1001, int64(1), "james")
	req.Header.Set("Authorization", token)

	ctx := context.WithValue(req.Context(), "roleID", float64(1))
	res := httptest.NewRecorder()
	handler := AuthAdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	handler.ServeHTTP(res, req.WithContext(ctx))
	t.Log(res.Code)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusOK, status)
	}
}

func TestAuthAdminMiddlewareFail(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		panic(err)
	}

	token := utils.CreateToken(1001, int64(2), "james")
	req.Header.Set("Authorization", token)

	ctx := context.WithValue(req.Context(), "roleID", float64(2))
	res := httptest.NewRecorder()
	handler := AuthAdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	handler.ServeHTTP(res, req.WithContext(ctx))
	t.Log(res.Code)
	if status := res.Code; status != http.StatusForbidden {
		t.Errorf("wrong status code: expected %v but got %v", http.StatusForbidden, status)
	}
}
