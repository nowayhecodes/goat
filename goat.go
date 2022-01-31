package goat

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

type HandlerFn func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFn
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

func (group *RouterGroup) Use(middlewares ...HandlerFn) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFn) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFn) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFn) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) PUT(pattern string, handler HandlerFn) {
	group.addRoute("PUT", pattern, handler)
}

func (group *RouterGroup) PATCH(pattern string, handler HandlerFn) {
	group.addRoute("PATCH", pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler HandlerFn) {
	group.addRoute("DELETE", pattern, handler)
}

func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFn) {
	group.addRoute("OPTIONS", pattern, handler)
}

func (group *RouterGroup) HEAD(pattern string, handler HandlerFn) {
	group.addRoute("HEAD", pattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFn {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(ctx *Context) {
		file := ctx.Param("filePath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}
