package goat

import "html/template"

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
