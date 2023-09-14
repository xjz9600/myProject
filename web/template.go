package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	Render(ctx context.Context, templateName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	Tpl *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, templateName string, data any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := g.Tpl.ExecuteTemplate(buffer, templateName, data)
	return buffer.Bytes(), err
}
