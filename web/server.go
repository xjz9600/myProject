package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 方法保证可以自己管理生命周期
	Start(add string) error
	AddRoute(method string, path string, handleFunc HandleFunc, ms ...Middleware)
}

var _ Server = &HttpServer{}

type ServerOption func(server *HttpServer)

func NewServer(opts ...ServerOption) *HttpServer {
	server := &HttpServer{
		router: NewRouter(),
	}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

func WithMiddleWare(mds ...Middleware) ServerOption {
	return func(server *HttpServer) {
		server.mds = mds
	}
}

type HttpServer struct {
	*router
	mds []Middleware
}

func (h *HttpServer) GET(path string, handleFunc HandleFunc, ms ...Middleware) {
	h.AddRoute(http.MethodGet, path, handleFunc, ms...)
}

func (h *HttpServer) POST(path string, handleFunc HandleFunc, ms ...Middleware) {
	h.AddRoute(http.MethodPost, path, handleFunc, ms...)
}

func (h *HttpServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 可以设置回调或者初始化一些数据
	return http.Serve(listener, h)
}

// ServeHTTP 自己路有在这里实现
func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := &Context{
		Response: writer,
		Req:      request,
	}
	root := h.Serve
	for i := len(h.mds) - 1; i >= 0; i-- {
		root = h.mds[i](root)
	}
	var finalStep Middleware = func(handleFunc HandleFunc) HandleFunc {
		return func(ctx *Context) {
			handleFunc(ctx)
			if ctx.RespStatusCode > 0 {
				ctx.Response.WriteHeader(ctx.RespStatusCode)
				ctx.Response.Write(ctx.RespData)
			}
		}
	}
	root = finalStep(root)
	root(context)
}

func (h *HttpServer) Serve(ctx *Context) {
	nodeInfo, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || nodeInfo.handle == nil {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.Response.Write([]byte("web: not found"))
		return
	}
	ctx.PathParams = nodeInfo.params
	ctx.MathRoute = nodeInfo.matchRoute
	root := nodeInfo.handle
	for i := len(nodeInfo.node.ms) - 1; i >= 0; i-- {
		root = nodeInfo.node.ms[i](root)
	}
	for i := len(nodeInfo.ms) - 1; i >= 0; i-- {
		root = nodeInfo.ms[i](root)
	}
	root(ctx)
}
