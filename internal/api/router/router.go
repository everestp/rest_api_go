package router

import (
	"net/http"

	"github.com/everestp/rest_api_go/internal/api/handlers"
)

// Router initializes all routes and returns a ServeMux.
func MainRouter() *http.ServeMux {
	tRouter := teacherRouter()
	sRouter := studentRouter()

	tRouter.Handle("/", sRouter)
	return  tRouter

}
