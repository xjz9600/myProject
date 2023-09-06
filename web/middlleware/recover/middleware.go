package recover

import "myProject/web"

type MiddlewareBuilder struct {
	statusCode int
	data       []byte
	log        func(ctx *web.Context)
}

func NewMiddlewareBuilder(statusCode int, data []byte, log func(ctx *web.Context)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		statusCode: statusCode,
		data:       data,
		log:        log,
	}
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(handleFunc web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.Response.WriteHeader(m.statusCode)
					ctx.Response.Write(m.data)
					m.log(ctx)
				}
			}()
			handleFunc(ctx)
		}
	}
}
