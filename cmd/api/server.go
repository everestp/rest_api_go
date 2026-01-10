package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	// Using an alias 'mw' for cleaner middleware calls
	mw "github.com/everestp/rest_api_go/internal/api/middlewares"
)

// User represents a system entity. 
// Future Reference: Marshalling (JSON Tags)
// Tags like `json:"id"` instruct the encoding/json package how to map 
// Go struct fields to JSON keys during API responses.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// rootHandler manages the base path.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Logic Block: Exact Path Validation
	// Since "/" is a catch-all in ServeMux, we filter out non-root paths.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Root API"))
}

// teacherHandler demonstrates dynamic path parsing.
func teacherHandler(w http.ResponseWriter, r *http.Request) {
	// Logic Block: Path Parameter Extraction
	// Trimming prefixes/suffixes manually is the standard library way 
	// to get IDs without a third-party router.
	path := strings.TrimPrefix(r.URL.Path, "/teacher/")
	userID := strings.TrimSuffix(path, "/")

	if userID != "" && userID != "teacher" {
		fmt.Fprintf(w, "Viewing Teacher Profile: %s", userID)
		return
	}
	w.Write([]byte("General Teacher Directory"))
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Executive Access Granted"))
}

func main() {
	// Configuration Constants
	const (
		certFile = "cert.pem"
		keyFile  = "key.pem"
		port     = ":3000"
	)

	// Logic Block: Routing Table
	// Future Reference: ServeMux 
	// A trailing slash (e.g., "/teacher/") acts as a subtree match,
	// allowing all paths starting with that prefix to be handled by this function.
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/teacher/", teacherHandler)
	mux.HandleFunc("/exec", execsHandler)

	// Logic Block: Middleware Onion Stacking
	// We wrap from inside-out. The last wrapper called is the first to execute.
	// Flow: ResponseTime -> SecurityHeaders -> Cors -> Mux
	var handler http.Handler = mux
	handler = mw.Cors(handler)
	handler = mw.SecurityHeaders(handler)
	handler = mw.ResponseTimeMiddleware(handler)

	// Future Reference: mTLS & TLS Hardening
	// tls.Config allows us to enforce modern security standards, 
	// such as disabling outdated TLS 1.0/1.1 versions.
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		PreferServerCipherSuites: true,
	}

	// Logic Block: Server Initialization
	// Customizing the server allows for timeouts, preventing "Slowloris" attacks.
	server := &http.Server{
		Addr:         port,
		Handler:      handler,
		TLSConfig:    tlsConfig,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Secure Server running on https://localhost%s\n", port)

	// Logic Block: Execution & Error Handling
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf("Critical Server Failure: %v", err)
	}
}