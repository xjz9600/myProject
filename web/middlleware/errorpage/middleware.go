package errorpage

import "myProject/web"

type MiddlewareBuilder struct {
	errorPage map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		errorPage: map[int][]byte{},
	}
}

func (m *MiddlewareBuilder) AddPage(code int, page []byte) {
	m.errorPage[code] = page
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(handleFunc web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			handleFunc(ctx)
			page, ok := m.errorPage[ctx.RespStatusCode]
			if ok {
				ctx.RespData = page
			}
		}
	}
}
