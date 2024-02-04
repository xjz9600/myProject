package main

import (
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"text/template"
)

//go:embed tpl.html
var genOrm string

func gen(w io.Writer, filePath string) error {
	fSet := token.NewFileSet()
	_, err := os.Open(filePath)
	f, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	fv := &FileVisitor{}
	ast.Walk(fv, f)
	file := fv.Get()
	tpl := template.New("gen-orm")
	tpl, err = tpl.Parse(genOrm)
	if err != nil {
		return err
	}
	return tpl.Execute(w, Data{
		File: file,
		Opts: []string{"LT", "GT", "EQ"},
	})
}

type Data struct {
	File
	Opts []string
}

func main() {

}
