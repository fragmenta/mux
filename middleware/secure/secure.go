// Package secure adds headers to protect against xss and reflection attacks and force use of https
package secure

import (
	"net/http"
)

// These package level variables should be called if required to set policies before the middleware is added

// ContentSecurityPolicy defaults to a strict policy disallowing iframes and scripts from any other origin save self (and Google Analytics for scripts)
var ContentSecurityPolicy = "frame-ancestors 'self'; style-src 'self'; script-src 'self' www.google-analytics.com"

// AddXHeaders determines whether the older headers X-XSS-Protection and X-Content-Type-Options are set - it defaults to true at present
var AddXHeaders = true

// Middleware adds some headers suitable for secure sites
func Middleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Add some headers for security

		// Allow no iframing - could also restrict scripts to this domain only (+GA?)
		w.Header().Set("Content-Security-Policy", ContentSecurityPolicy)

		// Allow only https connections for the next 30 days - this is ignored on http connections
		w.Header().Set("Strict-Transport-Security", "max-age=2592000")

		// Set ReferrerPolicy explicitly to send only the domain, not the path
		w.Header().Set("Referrer-Policy", "strict-origin")

		if AddXHeaders {
			// Ask browsers to block xss by default
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Don't allow browser sniffing for content types
			w.Header().Set("X-Content-Type-Options", "nosniff")
		}

		// Call the handler
		h(w, r)

	}
}
