package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// RequestData holds form data
func RequestData(r *http.Request, w http.ResponseWriter) map[string]interface{} {
	var data map[string]interface{}
	if len(r.Form) > 0 {
		for key, val := range r.Form {
			(data[key]) = strings.Join(val, "")
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Fprintln(w, err)
		}
		defer r.Body.Close()
		json.Unmarshal(body, &data)
	}
	return data
}

// ArrayRequestData holds array of form data
func ArrayRequestData(r *http.Request, w http.ResponseWriter) []map[string]interface{} {
	var data []map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	defer r.Body.Close()
	json.Unmarshal(body, &data)
	return data
}
