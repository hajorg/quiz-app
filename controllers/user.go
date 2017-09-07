package controllers

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"quiz-app/database"
	"quiz-app/utils"
	"quiz-app/validation"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserInput struct to replicate json request body
type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var dbName string

func init() {
	if flag.Lookup("test.v") == nil {
		dbName = "quiz"
	} else {
		dbName = "quiztest"
	}
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

	db := database.Connect(dbName)
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

	tokenString := utils.CreateToken(lastID, newUser["username"].(string))

	newUser["token"] = tokenString
	delete(newUser, "password")
	json.NewEncoder(w).Encode(newUser)
}

// Login login a registered user and give a token
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// holds form data
	user := map[string]interface{}{}
	// check if post data is urlencoded or json object
	if len(r.Form) > 0 {
		for key, val := range r.Form {
			user[key] = strings.Join(val, "")
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		defer r.Body.Close()
		json.Unmarshal(body, &user)
	}

	validationError := validation.Validator(w, user)
	if validationError == false {
		return
	}

	db := database.Connect(dbName)
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, username, password FROM user WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	var username string
	var password string

	err = stmt.QueryRow(user["username"]).Scan(&id, &username, &password)
	if err != nil {
		utils.UnauthorizedError(w, "Incorrect username and password combination")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user["password"].(string)))
	if err != nil {
		utils.UnauthorizedError(w, "Incorrect username and password combination")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	tokenString := utils.CreateToken(int64(id), username)
	user["token"] = tokenString

	json.NewEncoder(w).Encode(user)
}
