package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quiz-app/utils"
	"regexp"
	"strconv"
)

// Validator validates the request body(data) and writes to w if err
func Validator(w http.ResponseWriter, data map[string]interface{}) bool {
	validation := map[string](map[string]string){
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
		"title": {
			"required": "1",
			"min":      "2",
			"max":      "50",
		},
		"description": {
			"max": "255",
		},
		"name": {
			"required": "1",
			"min":      "2",
			"max":      "50",
		},
		"category_id": {
			"required": "1",
		},
	}

	// loop through the key of the outer map
	for attr, validate := range validation {
		// loop the value of the other loop with is also a map(inner map)
		for key, value := range validate {
			// loop through data to validate against
			for i, val := range data {
				if attr == i {
					switch key {
					case "required":
						if len(val.(string)) == 0 {
							return message(w, attr+" is required")
						}
					case "min":
						compare, _ := strconv.Atoi(value)
						if len(val.(string)) < compare {
							return message(w, fmt.Sprintf("%s should be atleast %d characters long", attr, compare))
						}
					case "max":
						compare, _ := strconv.Atoi(value)
						if len(val.(string)) > compare {
							return message(w, fmt.Sprintf("%s should not be more than %d characters", attr, compare))
						}
					case "pattern":
						match, err := regexp.MatchString(
							"(?i)[A-Za-z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Za-z]{2,4}",
							val.(string))
						if err != nil {
							panic(err)
						}
						if match == false {
							return message(w, fmt.Sprintf("%s is not a valid %s address", val, attr))
						}
					}
				}
			}
		}
	}
	return true
}

func message(w http.ResponseWriter, message string) bool {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	error := utils.Error{
		Error: message,
	}
	json.NewEncoder(w).Encode(error)
	return false
}
