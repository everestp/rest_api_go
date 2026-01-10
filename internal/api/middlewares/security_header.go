package middlewares

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Prevent clickjacking
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")

		// 2. Control DNS prefetching
		w.Header().Set("X-DNS-Prefetch-Control", "off")

		// 3. Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// 4. Basic XSS Protection (Note: Modern browsers rely more on CSP)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// 5. Referrer Policy: Only send origin for cross-site requests
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// 6. Permissions Policy: Disable sensitive features like camera/microchip by default
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// 7. Strict-Transport-Security (HSTS): Force HTTPS
		// Only use this if your site is served over HTTPS!
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// 8. Content Security Policy (CSP)
		// WARNING: Start with a 'report-only' or a permissive policy and tighten it.
		// This example allows scripts/styles from your own domain only.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self';")

		next.ServeHTTP(w, r)
	})
}


// Basic Middleware SKELETON
// func securityHeaders(next http.Handler) http.Handler{
//  return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//    next.ServeHTTP(w, r)
//  })

// }