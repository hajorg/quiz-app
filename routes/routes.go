package routes

import (
	"net/http"
	"quiz-app/controllers"
)

// Route holds required data to match an incoming request
type Route struct {
	Name    string
	Handler http.HandlerFunc
	Pattern string
	Method  string
}

type routes []Route

// Routers matches incoming request to the appropriate handlers
func Routers(w http.ResponseWriter, r *http.Request) {
	routes := routes{
		Route{
			Name:    "home",
			Handler: controllers.Index,
			Pattern: "/",
			Method:  "GET",
		},
		Route{
			Name:    "user",
			Handler: controllers.CreateUser,
			Pattern: "/user",
			Method:  "POST",
		},
		Route{
			Name:    "user",
			Handler: controllers.Login,
			Pattern: "/login",
			Method:  "POST",
		},
	}

	for _, route := range routes {
		if r.URL.Path == route.Pattern && r.Method == route.Method {
			route.Handler(w, r)
			return
		}
	}

	http.NotFound(w, r)
}
