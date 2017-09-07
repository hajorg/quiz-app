package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// UserInput struct to replicate json request body
type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// holds form data
	newUser := map[string]interface{}{}
	// check if post data is urlencoded or json object
	if len(r.Form) > 0 {
		for key, val := range r.Form {
			newUser[key] = strings.Join(val, "")
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Fprintln(w, err)
		}
		defer r.Body.Close()
		json.Unmarshal(body, &newUser)
	}

	validationError := validation.Validator(w, newUser)
	if validationError == false {
		return
	}

	db := database.Connect("quiz")
	defer db.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser["password"].(string)), 10)
	if err != nil {
		panic(err)
	}

	smt, err := db.Prepare("INSERT INTO user(id, username, email, password, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer smt.Close()

	res, err := smt.Exec(nil, newUser["username"], newUser["email"], hashedPassword, time.Now(), time.Now())
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	userID := lastID
	username := newUser["username"].(string)
	claims := MyCustomClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: 15000,
		},
	}
	mySigningKey := []byte("super-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println(err)
	}

	newUser["token"] = tokenString
	delete(newUser, "password")
	json.NewEncoder(w).Encode(newUser)
}

// MyCustomClaims customize jwt
type MyCustomClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
