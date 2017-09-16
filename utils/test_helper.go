package utils

import (
	"database/sql"
	"quiz-app/database"
)

// DbTestInit Creates a new db connection for test
func DbTestInit() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS quiztest")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE quiztest")
	if err != nil {
		panic(err)
	}
	database.CreateDatabase("quiztest")
	return db
}
