package server

import "net/http"

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Router provides minimal routing functionality with middleware support.
type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

// NewRouter returns a new Router instance.
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Use adds a middleware to the Router.
func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

// Handle registers a new route with the given pattern and handler.
func (r *Router) Handle(pattern string, handler http.Handler) {
	// Apply middleware chain in reverse order
	finalHandler := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}
	r.mux.Handle(pattern, finalHandler)
}

// ServeHTTP makes Router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
