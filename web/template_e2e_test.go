//go:build e2e

package web

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()
	server.GET("/", func(ctx *Context) {
		ctx.Response.Write([]byte("hello world"))
	})
	server.GET("/user/login", func(ctx *Context) {
		ctx.Response.Write([]byte("hello user login"))
	})
	server.POST("/user/haha", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Response.Write([]byte("hello user login"))
	})
	server.Start(":8083")
}

func TestMiddlewareServer(t *testing.T) {
	var testMiddleware Middleware = func(handleFunc HandleFunc) HandleFunc {
		return func(ctx *Context) {
			fmt.Println("这是公共的middleware开始")
			handleFunc(ctx)
			fmt.Println("这是公共的middleware结束")
		}
	}
	server := NewServer(WithMiddleWare(testMiddleware))
	server.GET("/a/b/c", func(ctx *Context) {
		ctx.RespJson(&TestStruct{
			Name: "middlewareByServer",
		})
	})
	var nodeMiddleware Middleware = func(handleFunc HandleFunc) HandleFunc {
		return func(ctx *Context) {
			fmt.Println("这是路由上的middleware开始")
			handleFunc(ctx)
			fmt.Println("这是路由上的middleware结束")
		}
	}
	server.GET("/a/b/d", func(ctx *Context) {
		ctx.RespJson(&TestStruct{
			Name: "middlewareByNode",
		})
	}, nodeMiddleware)
	var childMiddleware Middleware = func(handleFunc HandleFunc) HandleFunc {
		return func(ctx *Context) {
			fmt.Println("这是子节点上的middleware开始")
			handleFunc(ctx)
			fmt.Println("这是子节点上的middleware结束")
		}
	}
	server.GET("/a/b/d/m", func(ctx *Context) {
		ctx.RespJson(&TestStruct{
			Name: "middlewareByChildNode",
		})
	}, childMiddleware)
	server.Start(":8083")
}

type TestStruct struct {
	Name string `json:"name"`
}
