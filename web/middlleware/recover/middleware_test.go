//go:build e2e

package recover

import (
	"log"
	"myProject/web"
	"net/http"
	"testing"
)

func TestNewMiddleWare(t *testing.T) {
	panicRecover := NewMiddlewareBuilder(http.StatusOK, []byte(`
<html>
	<body>
		<h1>测试崩溃页面成功</h1>
	</body>
</html>`), func(ctx *web.Context) {
		log.Println("测试完成")
	})
	server := web.NewServer(web.WithMiddleWare(panicRecover.Build()))
	server.GET("/test/panic", func(ctx *web.Context) {
		panic("测试崩溃")
	})
	server.Start(":8085")
}
