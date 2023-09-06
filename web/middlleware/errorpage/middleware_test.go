//go:build e2e

package errorpage

import (
	"myProject/web"
	"net/http"
	"testing"
)

func TestNewMiddleWare(t *testing.T) {
	errPage := NewMiddlewareBuilder()
	errPage.AddPage(http.StatusOK, []byte(`
<html>
	<body>
		<h1>测试错误页面成功</h1>
	</body>
</html>`))
	server := web.NewServer(web.WithMiddleWare(errPage.Build()))
	server.GET("/test/errPage", func(ctx *web.Context) {
		ctx.RespJson(&User{
			Name: "xieJunZe",
		})
	})
	server.Start(":8085")
}

type User struct {
	Name string `json:"name"`
}
