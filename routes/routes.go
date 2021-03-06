package routes

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/quiz-app/controllers"
	"github.com/quiz-app/middlewares"
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

var baseURL = "/api/v1"

// Routers matches incoming request to the appropriate handlers
func Routers(w http.ResponseWriter, r *http.Request) {
	numReg, _ := regexp.Compile("\\d+")
	urlPaths := strings.Split(r.RequestURI, "/")
	lastPath := urlPaths[len(urlPaths)-1]

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
			Pattern: baseURL + "/user",
			Method:  "POST",
		},
		Route{
			Name:    "user",
			Handler: controllers.Login,
			Pattern: baseURL + "/login",
			Method:  "POST",
		},
		Route{
			Name:       "categories",
			Handler:    controllers.CreateCategory,
			Pattern:    baseURL + "/category",
			Method:     "POST",
			Middleware: middlewares.AuthMiddleware(middlewares.AuthAdminMiddleware(http.HandlerFunc(controllers.CreateCategory))),
		},
		Route{
			Name:       "subjects",
			Handler:    controllers.CreateSubject,
			Pattern:    baseURL + "/subject",
			Method:     "POST",
			Middleware: middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreateSubject)),
		},
		Route{
			Name:       "questions",
			Handler:    controllers.CreateQuestion,
			Pattern:    baseURL + "/question",
			Method:     "POST",
			Middleware: middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreateQuestion)),
		},
		Route{
			Name:       "options",
			Handler:    controllers.CreateOption,
			Pattern:    baseURL + "/option",
			Method:     "POST",
			Middleware: middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreateOption)),
		},
		Route{
			Name:    "subjects",
			Handler: controllers.GetSubjects,
			Pattern: baseURL + "/subject",
			Method:  "GET",
		},
		Route{
			Name:    "subject",
			Handler: controllers.GetSubject,
			Pattern: baseURL + "/subject/" + numReg.FindString(lastPath),
			Method:  "GET",
		},
		Route{
			Name:    "category",
			Handler: controllers.GetCategories,
			Pattern: baseURL + "/category",
			Method:  "GET",
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
