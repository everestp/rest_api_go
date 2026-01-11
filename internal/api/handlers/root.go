package handlers

import "net/http"

// rootHandler manages the base path.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to the Root API (HTTP Mode)"))
}
