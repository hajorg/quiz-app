package main

import (
	"fmt"
	"net/http"

	"github.com/quiz-app/routes"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// database.CreateDatabase("quiz")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", routes.Routers)
	fmt.Println(http.ListenAndServe(":8080", mux))
}
