package handlers

import (
	"fmt"
	"net/http"
)

type RouteGroup interface {
	Use(middlewares ...func(http.Handler) http.Handler)
	Handle(method, pattern string, handler http.Handler)
	HandleFunc(method, pattern string, handlerFn http.HandlerFunc)
}

type routeGroup struct {
	*http.ServeMux
	basePath    string
	middlewares []func(http.Handler) http.Handler
}

func (g *routeGroup) Use(middlewares ...func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *routeGroup) Handle(method, pattern string, handler http.Handler) {
	for _, middleware := range g.middlewares {
		handler = middleware(handler)
	}
	pattern = fmt.Sprintf("%s %s%s", method, g.basePath, pattern)
	g.ServeMux.Handle(pattern, handler)
}

func (g *routeGroup) HandleFunc(method, pattern string, handlerFn http.HandlerFunc) {
	g.Handle(method, pattern, handlerFn)
}
