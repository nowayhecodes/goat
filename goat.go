package goat

import (
	"html/template"
	"log"
)

type HandlersFn func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlersFn
		parent      *RouterGroup
		engine      *Engine
	}

	Engine struct {
		*RouterGroup
		router        *router
		groups        []*RouterGroup
		htmlTemplates *template.Template
		funcMap       template.FuncMap
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter} // <- TODO: define a router.go
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger, Recovery()) // <- TODO: define a recovery middleware
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlersFn) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlersFn) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlersFn) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlersFn) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) PUT(pattern string, handler HandlersFn) {
	group.addRoute("PUT", pattern, handler)
}

func (group *RouterGroup) PATCH(pattern string, handler HandlersFn) {
	group.addRoute("PATCH", pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler HandlersFn) {
	group.addRoute("DELETE", pattern, handler)
}

func (group *RouterGroup) OPTIONS(pattern string, handler HandlersFn) {
	group.addRoute("OPTIONS", pattern, handler)
}

func (group *RouterGroup) HEAD(pattern string, handler HandlersFn) {
	group.addRoute("HEAD", pattern, handler)
}
