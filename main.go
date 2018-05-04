package main

import (
	"fmt"
	"net/http"

	"github.com/quiz-app/routes"
	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// database.CreateDatabase("quiz")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", routes.Routers)
	handler := cors.Default().Handler(mux)
	fmt.Println(http.ListenAndServe(":8080", handler))
}
