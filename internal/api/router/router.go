package router

import (
	"net/http"

	"github.com/everestp/rest_api_go/internal/api/handlers"
)



func Router() *http.ServeMux{
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/teacher/", handlers.TeacherHandler)
	mux.HandleFunc("/students/", handlers.StudentHandler)
	mux.HandleFunc("/execes/", handlers.ExcesHandler)

}
