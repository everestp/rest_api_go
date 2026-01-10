package middlewares

import (

	"net/http"
)


func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Set the allowed origin
		// In production, replace "*" with specific domains like "https://myapp.com"
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 2. Specify which methods are allowed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 3. Specify which headers the client can send
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 4. Handle "Preflight" requests
		// Browsers send an OPTIONS request before POST/PUT/DELETE
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
