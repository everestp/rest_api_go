package router

import (
	"net/http"

	"github.com/everestp/rest_api_go/internal/api/handlers"
)

// Router initializes all routes and returns a ServeMux.
func Router() *http.ServeMux {
	mux := http.NewServeMux()

	// Root
	mux.HandleFunc("/", handlers.RootHandler)

	// Teacher endpoints
	mux.HandleFunc("/teacher/", handlers.TeacherHandler)

	// Student endpoints
	mux.HandleFunc("/students/", handlers.StudentHandler)

	// Exercise endpoints
	mux.HandleFunc("/execes/", handlers.ExcesHandler)

	return mux
}
