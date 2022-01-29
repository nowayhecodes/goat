package goat

import (
	"log"
	"time"
)

func Logger() FuncHandler {
	return func(ctx *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
