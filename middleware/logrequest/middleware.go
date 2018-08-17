package logrequest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fragmenta/mux/log"
)

// TargetResponseTime sets the threshold for colorisation of response times
var TargetResponseTime = 50 * time.Millisecond

// Middleware logs after each request to record the method, the url, the status code and the response time
// e.g. GET / -> status 200 in 31.932146ms
func Middleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Store the time prior to handling
		start := time.Now()

		// Wrap the response writer to record code
		// Ideally we'd instead take mux.HandlerFunc
		cw := newCodeResponseWriter(w)

		// Run the handler with our recording response writer
		h(cw, r)

		// Calculate method, url, code, response time
		method := r.Method
		url := r.URL.Path
		duration := time.Now().UTC().Sub(start)
		code := cw.StatusCode

		// Skip logging assets, favicon
		if strings.HasPrefix(url, "/assets") || strings.HasPrefix(url, "/favicon.ico") {
			return
		}

		// Pretty print to the standard loggers colorized
		logWithColor(method, url, code, duration)

		// Log the values to any value loggers (for export to monitoring services)
		values := map[string]interface{}{
			"method": r.Method,
			"url":    r.URL.Path,
			"code":   code,
			"time":   duration,
		}
		log.Values(values)
	}

}

// codeResponseWriter defines a responseWriter which stores the status code
type codeResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader stores the code before writing
func (cw *codeResponseWriter) WriteHeader(code int) {
	cw.StatusCode = code
	cw.ResponseWriter.WriteHeader(code)
}

// newCodeResponseWriter initialises a codeResponseWriter
func newCodeResponseWriter(w http.ResponseWriter) *codeResponseWriter {
	return &codeResponseWriter{w, http.StatusOK}
}

// Format a string by wrapping in a given color code
func applyColor(f, s string) string {
	return f + s + log.ColorNone
}

// logWithColor formats the log string with color depending on the arguments
func logWithColor(method string, url string, code int, duration time.Duration) {

	// Start with all green, colorise output depending on values
	m := log.ColorGreen
	c := log.ColorGreen
	d := log.ColorGreen

	// Only GET is green
	if method != http.MethodGet {
		m = log.ColorAmber
	}

	// Only 200 is green
	if code != http.StatusOK {
		c = log.ColorRed
	}

	// Only under TargetResponseTime is green
	if duration > TargetResponseTime {
		d = log.ColorRed
	}

	// Generate a format string using colors to wrap formats for values
	// The equivalent of the plain format "%s %s -> %d in %s"
	format := fmt.Sprintf("%s %%s %s %s in %s", applyColor(m, "%s"), applyColor(log.ColorCyan, "->"), applyColor(c, "%d"), applyColor(d, "%s"))

	// Print to the log with this colorised format
	log.Printf(format, method, url, code, duration)
}
