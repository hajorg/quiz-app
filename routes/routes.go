package routes

import (
	"net/http"
	"quiz-app/controllers"
	"quiz-app/middlewares"
)

// Route holds required data to match an incoming request
type Route struct {
	Name       string
	Handler    http.HandlerFunc
	Pattern    string
	Method     string
	Middleware http.Handler
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
		Route{
			Name:       "categories",
			Handler:    controllers.CreateCategory,
			Pattern:    "/category",
			Method:     "POST",
			Middleware: middlewares.AuthMiddleware(middlewares.AuthAdminMiddleware(http.HandlerFunc(controllers.CreateCategory))),
		},
	}

	for _, route := range routes {
		if r.URL.Path == route.Pattern && r.Method == route.Method {
			if route.Middleware != nil {
				route.Middleware.ServeHTTP(w, r)
			} else {
				route.Handler(w, r)
			}
			return
		}
	}

	http.NotFound(w, r)
}
