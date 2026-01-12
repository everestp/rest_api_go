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
	mux.HandleFunc("GET /teacher/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teacher/", handlers.AddTeacherHandler)
	mux.HandleFunc("PATCH /teacher/", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teacher/", handlers.DeleteTeachersHandler)

	mux.HandleFunc("GET /teacher/{id}", handlers.GetTeacherHandler)
	mux.HandleFunc("PUT /teacher/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teacher/{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teacher/{id}", handlers.DeleteOneTeacherHandler)

	// Student endpoints
	mux.HandleFunc("/students/", handlers.StudentHandler)

	// Exercise endpoints
	mux.HandleFunc("/execes/", handlers.ExcesHandler)

	return mux
}
