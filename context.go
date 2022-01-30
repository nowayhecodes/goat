package goat

import (
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	handlers   []FuncHandler
	index      int
	engine     *Engine
}

func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: writer,
		index:  -1,
	}
}

func (ctx *Context) Next() {
	ctx.index++
	h := len(ctx.handlers)

	for ; ctx.index < h; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}
