package main

import (
	"net/http"
	"quiz-app/database"
	"quiz-app/routes"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	database.CreateDatabase("quiz")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", routes.Routers)
	http.ListenAndServe(":8080", mux)
}
