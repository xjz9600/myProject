package file

import (
	"myProject/web"
	"net/http"
	"path/filepath"
	"strings"
)

type downloader struct {
	dst string
}

func NewDownLoader(dst string) *downloader {
	return &downloader{dst: dst}
}

func (d *downloader) Build() web.HandleFunc {
	return func(ctx *web.Context) {
		fileName, err := ctx.QueryValue("file").AsString()
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte(err.Error())
			return
		}
		path := filepath.Join(d.dst, fileName)
		finalPath, err := filepath.Abs(path)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte(err.Error())
			return
		}
		basePath, err := filepath.Abs(d.dst)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte(err.Error())
			return
		}
		if !strings.Contains(finalPath, basePath) {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("路径不合法")
			return
		}
		fn := filepath.Base(finalPath)
		header := ctx.Response.Header()
		// 必须设置
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		//设置是否使用缓存
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Response, ctx.Req, finalPath)
	}
}
