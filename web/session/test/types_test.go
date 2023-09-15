package test

import (
	"fmt"
	"myProject/web"
	"myProject/web/session"
	"myProject/web/session/memory"
	"myProject/web/session/propagator"
	"net/http"
	"testing"
)

func TestTypes(t *testing.T) {
	var manager = session.Manager{
		Store:         memory.NewStore(),
		Propagator:    propagator.NewPropagator(),
		CtxSessionKey: "sessionCache",
	}
	server := web.NewServer(web.WithMiddleWare(func(handleFunc web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			if ctx.Req.URL.Path == "/login" {
				handleFunc(ctx)
				return
			}
			err := manager.RefreshSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusInternalServerError
				ctx.RespData = []byte(err.Error())
				return
			}
			handleFunc(ctx)
		}
	}))
	server.POST("/login", func(ctx *web.Context) {
		session, err := manager.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte(err.Error())
			return
		}
		session.Set(ctx.Req.Context(), "nickname", "xieJunZe")
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功")
	})
	server.GET("/logout", func(ctx *web.Context) {
		err := manager.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte(err.Error())
			return
		}
	})
	server.GET("/user", func(ctx *web.Context) {
		session, err := manager.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte(err.Error())
			return
		}
		val, err := session.Get(ctx.Req.Context(), "nickname")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte(err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte(fmt.Sprintf("欢迎：%v您的到来", val))
	})
	server.Start(":8083")
}
