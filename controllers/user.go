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

	validationError := validation.Validator(w, newUser, map[string](map[string]string){
		"username": {
			"required": "1",
			"max":      "20",
		},
		"email": {
			"required": "1",
			"pattern":  "1",
		},
		"password": {
			"required": "1",
			"min":      "6",
		},
	})
	if validationError == false {
		return
	}

	db := database.Connect(dbName)
	defer db.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser["password"].(string)), 10)
	if err != nil {
		panic(err)
	}
	newUser["password"] = hashedPassword
	newUser["created_at"] = time.Now()
	newUser["updated_at"] = time.Now()

	lastID, err := database.Insert("user", newUser)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	stmt, err := db.Prepare("SELECT id, username, email, role_id FROM user WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	var username string
	var email string
	var roleID int

	err = stmt.QueryRow(lastID).Scan(&id, &username, &email, &roleID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	tokenString := utils.CreateToken(lastID, int64(roleID), newUser["username"].(string))

	newUser["token"] = tokenString
	delete(newUser, "password")
	delete(newUser, "created_at")
	delete(newUser, "updated_at")
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

	validationError := validation.Validator(w, user, map[string](map[string]string){
		"username": {
			"required": "1",
			"max":      "20",
		},
		"password": {
			"required": "1",
			"min":      "6",
		},
	})
	if validationError == false {
		return
	}

	db := database.Connect(dbName)
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, username, password, role_id FROM user WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	var username string
	var password string
	var roleID int

	err = stmt.QueryRow(user["username"]).Scan(&id, &username, &password, &roleID)
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

	tokenString := utils.CreateToken(int64(id), int64(roleID), username)
	user["token"] = tokenString
	delete(user, "password")

	json.NewEncoder(w).Encode(user)
}
