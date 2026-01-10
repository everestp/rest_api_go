package middlewares

import (
    "fmt"
    "net/http"
    "strings"
)

// HPPOptions defines the configuration for the HTTP Parameter Pollution middleware.
type HPPOptions struct {
    CheckQuery                  bool     // Logic Block: Enable/Disable scanning of URL query strings
    CheckBody                   bool     // Logic Block: Enable/Disable scanning of POST body (form-data)
    CheckBodyOnlyForContentType string   // Logic Block: Filter by Content-Type (e.g., "application/x-www-form-urlencoded")
    Whitelist                   []string // Logic Block: Parameters allowed to pass through without being deleted
}

// Hpp initializes the middleware with specific security constraints.
// Future Reference: Parameter Pollution
// HPP occurs when multiple parameters with the same name are sent (e.g., ?id=1&id=2). 
// Different servers handle this differently (first value, last value, or array). 
// This middleware enforces a predictable "Single Value" behavior.
func Hpp(options HPPOptions) func(http.Handler) http.Handler {
    fmt.Println("HPP Middleware Initialized...")
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            
            // Logic Block: Body Parameter Filtering
            // We only check the body if it's a POST request and matches the expected Content-Type.
            if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
                filterBodyParams(r, options.Whitelist)
            }

            // Logic Block: Query Parameter Filtering
            if options.CheckQuery && r.URL.Query() != nil {
                filterQueryParams(r, options.Whitelist)
            }

            next.ServeHTTP(w, r)
        })
    }
}

// isCorrectContentType validates the header to prevent unnecessary parsing of JSON or XML as Forms.
func isCorrectContentType(r *http.Request, contentType string) bool {
    return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

// filterBodyParams parses the form and reduces multi-value keys to a single value.
func filterBodyParams(r *http.Request, whitelist []string) {
    // Future Reference: ParseForm
    // ParseForm populates r.Form. It handles both URL query parameters and POST form data.
    err := r.ParseForm()
    if err != nil {
		// In production, consider sending an http.Error instead of just printing
        fmt.Println("Error parsing form:", err)
        return
    }

    for k, v := range r.Form {
        // Logic Block: Flattening
        // If an attacker sends ?user=admin&user=guest, we force it to just "admin".
        if len(v) > 1 {
            r.Form.Set(k, v[0]) 
        }
        
        // Logic Block: Whitelisting
        // If the parameter key isn't in our allowed list, we remove it entirely.
        if len(whitelist) > 0 && !isWhiteListed(k, whitelist) {
            delete(r.Form, k)
        }
    }
}

// filterQueryParams flattens and cleans the URL string.
func filterQueryParams(r *http.Request, whitelist []string) {
    query := r.URL.Query()

    for k, v := range query {
        if len(v) > 1 {
            // Logic Block: Last-Value Win
            // Your implementation currently chooses the LAST value for queries.
            query.Set(k, v[len(v)-1]) 
        }
        
        if len(whitelist) > 0 && !isWhiteListed(k, whitelist) {
            query.Del(k)
        }
    }
    // Future Reference: URL Re-encoding
    // After modifying the map, we must re-encode the string back into the Request object.
    r.URL.RawQuery = query.Encode()
}

// isWhiteListed performs a simple linear search to validate the parameter key.
func isWhiteListed(param string, whitelist []string) bool {
    for _, v := range whitelist {
        if param == v {
            return true
        }
    }
    return false
}