package accesslog

import (
	"encoding/json"
	"myProject/web"
)

type MiddleWareBuilder struct {
	logFun func(msg string)
}

func NewMiddleWare(logFun func(msg string)) *MiddleWareBuilder {
	return &MiddleWareBuilder{
		logFun: logFun,
	}
}

type accessLog struct {
	Host       string `json:"host,omitempty"`
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}

func (m *MiddleWareBuilder) Build() web.Middleware {
	return func(handleFunc web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
					Route:      ctx.MathRoute,
				}
				data, _ := json.Marshal(l)
				m.logFun(string(data))
			}()
			handleFunc(ctx)
		}
	}
}
