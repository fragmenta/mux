// Package secure adds headers to protect against xss and reflection attacks and force use of https
package secure

import (
	"net/http"
)

// Middleware adds some headers suitable for secure sites
func Middleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Add some headers for security

		// Allow no iframing - could also restrict scripts to this domain only (+GA?)
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'; style-src 'self'; script-src 'self' www.google-analytics.com")

		// Allow only https connections for the next 30 days - this is ignored on http connections
		w.Header().Set("Strict-Transport-Security", "max-age=2592000")

		// Set ReferrerPolicy explicitly to send only the domain, not the path
		w.Header().Set("ReferrerPolicy", "strict-origin")

		// Ask browsers to block xss by default
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Don't allow browser sniffing for content types
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Call the handler
		h(w, r)

	}
}
