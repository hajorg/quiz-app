package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/quiz-app/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// AuthMiddleware takes the token provided by the user and validates it
func AuthMiddleware(h http.Handler) http.Handler {
	godotenv.Load()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" || len(strings.Split(tokenString, ".")) != 3 {
			utils.UnauthorizedError(w, "No Token provided")
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.UnauthorizedError(w, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "roleID", claims["roleId"])

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthAdminMiddleware checks if the user is an admin
func AuthAdminMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleID := r.Context().Value("roleID").(float64)

		if roleID != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(403)
			w.Write([]byte(`{"error": "You are not permitted to perform this action"}`))
			return
		}
		h.ServeHTTP(w, r)
	})
}
