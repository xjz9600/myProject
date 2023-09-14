//go:build e2e

package file

import (
	"github.com/stretchr/testify/assert"
	"html/template"
	"myProject/web"
	"testing"
)

func TestUploader(t *testing.T) {
	tpl, err := template.ParseGlob("../testdata/tpls/*.gohtml")
	assert.NoError(t, err)
	tmp := &web.GoTemplateEngine{
		Tpl: tpl,
	}
	server := web.NewServer(web.WithTplEngine(tmp))
	server.GET("/uploader", func(ctx *web.Context) {
		ctx.Render("upload.gohtml", nil)
	})
	server.POST("/uploader", NewFileUploader().Handle)
	server.Start(":8083")
}
