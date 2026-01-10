package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	mw "github.com/everestp/rest_api_go/internal/api/middlewares"
)

// rootHandler manages the base path.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to the Root API (HTTP Mode)"))
}

// teacherHandler demonstrates dynamic path parsing.
func teacherHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teacher/")
	userID := strings.TrimSuffix(path, "/")

	if userID != "" && userID != "teacher" {
		fmt.Fprintf(w, "Viewing Teacher Profile: %s", userID)
		return
	}
	w.Write([]byte("General Teacher Directory"))
}

func main() {
	// Logic Block: Configuration
	// We no longer need certFile or keyFile constants.
	const port = ":3000"

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/teacher/", teacherHandler)
   rl := mw.NewRateLimiter(5, time.Minute)
     hppOptions := mw.HPPOptions{
		CheckQuery: true,
		CheckBody: true,
		CheckBodyOnlyForContentType: "application/x-www-from-urlencode",
		Whitelist: []string{"name"},
	 }
	 fmt.Println(rl ,hppOptions)
	// Logic Block: Middleware Onion
	// The order remains the same: Timing -> Compression -> Security -> CORS -> App
	// secureMux1 := applyMiddlewares(mux, mw.Hpp(hppOptions) ,mw.Compression , mw.SecurityHeaders , mw.ResponseTimeMiddleware , rl.Middleware ,mw.Cors)
   secureMux := mw.SecurityHeaders(mux)
	// Logic Block: Server Initialization
	// We removed the TLSConfig field.
	server := &http.Server{
		Addr:         port,
		Handler:      secureMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Future Reference: HTTP vs HTTPS
	// Without TLS, data is sent in 'Plaintext'. Anyone on the network 
	// can see the traffic. This is fine for local dev but never for production.
	fmt.Printf("🚀 Server running on http://localhost%s\n", port)

	// Logic Block: Standard Execution
	// We use ListenAndServe() instead of ListenAndServeTLS()
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Critical Server Failure: %v", err)
	}
}

//Middleware is the  fuction that wraps an http.Handler with additional functionality
type Middleware func(http.Handler) http.Handler


func ApplyMiddlewares(handler http.Handler , middlewares ...Middleware) http.Handler{
 for _ , middleware := range middlewares{
	handler = middleware(handler)
 }
 return  handler
}