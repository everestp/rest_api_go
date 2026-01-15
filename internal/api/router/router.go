package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	tRouter := teacherRouter()
	sRouter := studentRouter()
	sRouter.Handle("/", execsRouter())
	tRouter.Handle("/", sRouter)
	return tRouter
}