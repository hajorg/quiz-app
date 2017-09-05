package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
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
	fmt.Println(r, "here")
	var user UserInput
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	defer r.Body.Close()
	var newUser map[string]interface{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &newUser)
	fmt.Println(newUser)
	validationError := validation.Validator(w, newUser)
	if validationError == false {
		return
	}

	db := database.Connect()
	defer db.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		panic(err)
	}

	smt, err := db.Prepare("INSERT INTO user(id, username, email, password, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer smt.Close()

	res, err := smt.Exec(nil, user.Username, user.Email, hashedPassword, time.Now(), time.Now())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		error := utils.Error{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(error)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		error := utils.Error{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(error)
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

	fmt.Println(tokenString)
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
