package http

import (
	"net/http"
	"slices"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	mux             *http.ServeMux
	middlewares     []Middleware
	prefix          string
	notFoundHandler http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: make([]Middleware, 0),
		prefix:      "",
	}
}

// middlewares and composition

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) With(m Middleware) *Router {
	clone := r.clone()
	clone.middlewares = append(clone.middlewares, m)

	return clone
}

func (r *Router) Group(fn func(*Router)) *Router {
	clone := r.clone()
	fn(clone)

	return clone
}

func (r *Router) Route(pattern string, fn func(*Router)) *Router {
	clone := r.clone()

	clone.prefix = strings.TrimSuffix(r.prefix, "/") + "/" + strings.TrimPrefix(pattern, "/")
	fn(clone)

	return clone
}

func (r *Router) Mount(pattern string, handler http.Handler) {
	mountPath := strings.TrimPrefix(r.prefix+pattern, "/") + "/"

	r.mux.Handle(mountPath, handler)
}

// helpers

func (r *Router) clone() *Router {
	mCopy := make([]Middleware, len(r.middlewares))
	copy(mCopy, r.middlewares)

	return &Router{
		mux:             r.mux,
		middlewares:     mCopy,
		prefix:          r.prefix,
		notFoundHandler: r.notFoundHandler,
	}
}

func (r *Router) handle(method, path string, handler http.HandlerFunc) {
	var finalHandler http.Handler = handler

	for _, middleware := range slices.Backward(r.middlewares) {
		finalHandler = middleware(finalHandler)
	}

	fullPath := r.prefix + path
	if fullPath == "" {
		fullPath = "/"
	}

	r.mux.Handle(method+" "+path, finalHandler)
}

// http methods

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.handle(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.handle(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.handle(http.MethodPut, path, handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.handle(http.MethodDelete, path, handler)
}

func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.handle(http.MethodPatch, path, handler)
}

func (r *Router) Head(path string, handler http.HandlerFunc) {
	r.handle(http.MethodHead, path, handler)
}

func (r *Router) Options(path string, handler http.HandlerFunc) {
	r.handle(http.MethodOptions, path, handler)
}

func (r *Router) Connect(path string, handler http.HandlerFunc) {
	r.handle(http.MethodConnect, path, handler)
}

func (r *Router) Trace(path string, handler http.HandlerFunc) {
	r.handle(http.MethodTrace, path, handler)
}

// custom error handling

func (r *Router) NotFound(handler http.HandlerFunc) {
	r.notFoundHandler = handler
}

// server

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, pattern := r.mux.Handler(req)

	if pattern == "" && r.notFoundHandler != nil {
		r.notFoundHandler(w, req)
		return
	}

	handler.ServeHTTP(w, req)
}
