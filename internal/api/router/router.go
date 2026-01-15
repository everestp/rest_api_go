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
	mux.HandleFunc("GET /teacher", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teacher", handlers.AddTeacherHandler)
	mux.HandleFunc("PATCH /teacher", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teacher", handlers.DeleteTeachersHandler)

	mux.HandleFunc("GET /teacher/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /teacher/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teacher/{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teacher/{id}", handlers.DeleteOneTeacherHandler)

		mux.HandleFunc("GET /teacher/{id}/students", handlers.GetStudentByTeacherID)
		mux.HandleFunc("GET /teacher/{id}/studentscount", handlers.GetOneTeacherHandler)

	// Student endpoints
	// mux.HandleFunc("/students/", handlers.StudentHandler)
	mux.HandleFunc("GET /students", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students", handlers.AddStudentHandler)
	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", handlers.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchOneStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteOneStudentHandler)

	// Exercise endpoints
	mux.HandleFunc("/execes/", handlers.ExcesHandler)

	return mux
}
