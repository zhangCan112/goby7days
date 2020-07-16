package gee

import (
	"log"
	"net/http"
)

type router struct {
	handers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
