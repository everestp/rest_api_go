package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Compression middleware detects if a client supports Gzip and compresses the response.
// Future Reference: Bandwidth Optimization
// Compression reduces the size of the payload (often by 70-80% for JSON), 
// improving load times and reducing data costs.
func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logic Block: Content-Negotiation
		// We check the "Accept-Encoding" header to see if the client can handle Gzip.
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Logic Block: Gzip Initialization
		// We set the header first to tell the browser "the body you are about to get is compressed."
		w.Header().Set("Content-Encoding", "gzip")
		
		gz := gzip.NewWriter(w)
		// Future Reference: Resource Cleanup
		// We MUST close the gzip writer to flush the remaining bytes to the response.
		defer gz.Close()

		// Logic Block: Dependency Injection
		// We wrap the original ResponseWriter with our custom Gzip writer.
		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		// Continue the chain with the wrapped writer
		next.ServeHTTP(gzw, r)
	})
}

// gzipResponseWriter intercepts the Write calls to pass data through the Gzip engine.
// Future Reference: Interface Promotion
// Because we embed http.ResponseWriter, this struct satisfies the interface 
// while overriding the Write method.
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write intercepts the standard Write call.
// Logic Block: Stream Redirection
func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	// Instead of writing to the network directly, we write to the Gzip compressor.
	return g.Writer.Write(b)
}