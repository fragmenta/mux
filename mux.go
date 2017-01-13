package mux

import (
	"net/http"
)

// HandlerFunc defines a std net/http HandlerFunc, but which returns an error.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorHandlerFunc defines a HandlerFunc which accepts an error and displays it.
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

// Middleware is a handler that wraps another handler
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Route defines the interface routes are expected to conform to.
type Route interface {
	MatchMethod(string) bool
	MatchMaybe(string) bool
	Match(string) bool
	String() string
	Handler() HandlerFunc
	Parse(string) map[string]string

	Get() Route
	Post() Route
	Put() Route
	Delete() Route
	Methods(...string) Route
}

// Mux handles http requests by selecting a handler
// and passing the request to it.
// Routes are evaluated in the order they were added.
// Before the request reaches the handler
// it is passed through the middleware chain.
type Mux struct {
	routes       []Route
	handlerFuncs []Middleware

	// See httptrace for best way to instrument
	ErrorHandler ErrorHandlerFunc
	FileHandler  HandlerFunc
}

// New returns a new mux
func New() *Mux {
	m := &Mux{
		FileHandler:  fileHandler,
		ErrorHandler: errHandler,
	}

	return m
}

// ServeHTTP implements net/http.Handler.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := m.RouteRequest
	for _, mh := range m.handlerFuncs {
		h = mh(h)
	}
	h(w, r)
}

// RouteRequest is the final endpoint of all requests
func (m *Mux) RouteRequest(w http.ResponseWriter, r *http.Request) {
	// Match a route
	route := m.Match(r)
	if route == nil {
		err := m.FileHandler(w, r)
		if err != nil {
			m.ErrorHandler(w, r, err)
		}
		return
	}

	// Execute the route
	err := route.Handler()(w, r)
	if err != nil {
		m.ErrorHandler(w, r, err)
	}

}

// Match finds the route (if any) which matches this request
func (m *Mux) Match(r *http.Request) Route {
	// Handle nil request
	if r == nil {
		return nil
	}

	// Routes are checked in order against the request path
	for _, route := range m.routes {
		// Test with probabalistic match
		if route.MatchMaybe(r.URL.Path) {
			// Test on method
			if route.MatchMethod(r.Method) {
				// Test exact match (may be expensive regexp)
				if route.Match(r.URL.Path) {
					return route
				}
			}

		}
	}

	return nil
}

// AddMiddleware adds a middleware function, this should be done before
// starting the server as it remakes our chain of middleware.
func (m *Mux) AddMiddleware(middleware Middleware) {
	// Prepend to our array of middleware
	m.handlerFuncs = append([]Middleware{middleware}, m.handlerFuncs...)
}

// Add adds a route for this request with the default methods (GET/HEAD)
// Route is returned so that method functions can be chained
func (m *Mux) Add(pattern string, handler HandlerFunc) Route {
	route, err := NewRoute(pattern, handler)
	if err != nil {
		// errors should be rare, but log them to stdout for debug
		println("mux: error parsing route:%s", pattern)
	}

	m.routes = append(m.routes, route)
	return route
}
