package rest

import (
	"log"
	"net/http"
	"regexp"
	"time"
)

type route struct {
	method  string
	path    *regexp.Regexp
	handler func(w http.ResponseWriter, r *http.Request)
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

type Router struct {
	routes []*route
}

func (h *Router) GET(path string, handler handlerFunc) error {
	return h.handleFunc("GET", path, handler)
}

func (h *Router) PUT(path string, handler handlerFunc) error {
	return h.handleFunc("PUT", path, handler)
}

func (h *Router) POST(path string, handler handlerFunc) error {
	return h.handleFunc("POST", path, handler)
}

func (h *Router) DELETE(path string, handler handlerFunc) error {
	return h.handleFunc("DELETE", path, handler)
}

func (h *Router) handleFunc(method string, path string, handler handlerFunc) error {
	pathRegExp, err := regexp.Compile(path)
	if err != nil {
		return err
	}
	h.routes = append(h.routes, &route{method: method, path: pathRegExp, handler: handler})
	return nil
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if r.Method == route.method && route.path.MatchString(r.URL.Path) {
			start := time.Now()
			route.handler(w, r)
			took := time.Since(start)
			log.Printf("[%v] %v %v [%s]", r.RemoteAddr, r.Method, r.URL.Path, took)
			return
		}
	}
	http.NotFound(w, r)
}
