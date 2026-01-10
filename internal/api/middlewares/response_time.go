package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

// ResponseTimeMiddleware captures the latency of an HTTP request and logs metadata.
// It uses the "Decorator" pattern to wrap the handler and observe its behavior.
func ResponseTimeMiddleware(next http.Handler) http.Handler {
	// Logic Block: Middleware Initialization
	// This code runs only ONCE when the server starts and the middleware chain is built.
	fmt.Println("Response Time Middleware Initialized...")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logic Block: Request Entry
		// This inner function runs for EVERY incoming request.
		fmt.Println("Response Time Middleware: Intercepting Request...")
		start := time.Now()

		// Logic Block: Dependency Injection / Wrapping
		// We substitute the standard ResponseWriter with our custom 'responseWriter' spy.
		// Future Reference: Interface Satisfaction
		// In Go, any struct that implements the methods Header(), Write(), and WriteHeader() 
		// satisfies the http.ResponseWriter interface.
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Logic Block: Hand-off to the Onion Layers
		// next.ServeHTTP passes the request down to the next middleware or the final handler.
		// Control "pauses" here until the rest of the application is finished.
		next.ServeHTTP(wrappedWriter, r)

		// Logic Block: Post-Execution Metrics
		// Once the handler returns, we calculate how long it took to process the entire request.
		duration := time.Since(start)

		// Future Reference: Header Mutability
		// Note: Setting headers after next.ServeHTTP() only works if the handler hasn't 
		// already sent (flushed) the response body to the client.
		w.Header().Set("X-Response-Time", duration.String())

		// Logic Block: Structured Logging
		// We use the 'status' captured by our wrappedWriter to log the result.
		fmt.Printf("Method: %s, URL: %s, Status: %d, Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Response Time Middleware: Request Cycle Complete")
	})
}

// Future Reference: Struct Embedding (Promotion)
// By embedding http.ResponseWriter, 'responseWriter' inherits all methods 
// of the interface automatically. We only "override" the ones we need to spy on.
type responseWriter struct {
	http.ResponseWriter // Anonymous embedding
	status int          // Field to store the intercepted status code
}

// WriteHeader captures the status code sent by the handler.
// Logic Block: Interception Logic
func (rw *responseWriter) WriteHeader(code int) {
	// Store the code so the middleware can read it later
	rw.status = code
	// Call the original ResponseWriter's method to actually send the code to the client
	rw.ResponseWriter.WriteHeader(code)
}