package database

import (
	"database/sql"
	"log"
)

// Connect connects to the local database
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
