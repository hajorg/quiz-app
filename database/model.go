package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
)

// Connect connects to the local database
var dbName string

func init() {
	if flag.Lookup("test.v") == nil {
		dbName = "quiz"
	} else {
		dbName = "quiztest"
	}
}

func Connect(name string) *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/"+name)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Insert(table string, data map[string]interface{}) (int64, error) {
	var newData []map[string]interface{}

	for key, val := range data {
		newData = append(newData, map[string]interface{}{key: val})
	}

	db := Connect(dbName)
	defer db.Close()

	placeholder := "" // used as prepared statement binding i.e ?
	query := "INSERT INTO " + fmt.Sprint(table) + "("

	content := []interface{}{}
	for idx, val := range newData {
		for key := range val {
			query += ("" + fmt.Sprint(key))
			placeholder += "?"
			content = append(content, val[key])
			if idx != len(newData)-1 {
				query += ", "
				placeholder += ", "
			}
		}
	}

	query += ") VALUES(" + placeholder + ")"
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(content...)

	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return lastID, nil
}
