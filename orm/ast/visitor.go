package main

import (
	"go/ast"
)

type FileVisitor struct {
	fileInfo *FileVisitorInfo
}

type File struct {
	Package string
	Imports []string
	Types   []Type
}

type Type struct {
	Name   string
	Fields []Field
}

func (f *FileVisitor) Get() File {
	types := make([]Type, 0, len(f.fileInfo.types))
	for _, typ := range f.fileInfo.types {
		types = append(types, Type{
			Name:   typ.Name,
			Fields: typ.Field,
		})
	}
	return File{
		Package: f.fileInfo.PackageName,
		Imports: f.fileInfo.Imports,
		Types:   types,
	}
}

func (f *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	file, ok := node.(*ast.File)
	if ok {
		f.fileInfo = &FileVisitorInfo{
			PackageName: file.Name.String(),
		}
		return f.fileInfo
	}
	return f
}

type FileVisitorInfo struct {
	PackageName string
	Imports     []string
	types       []*TypeVisitor
}

func (f *FileVisitorInfo) Visit(node ast.Node) (w ast.Visitor) {
	switch expr := node.(type) {
	case *ast.ImportSpec:
		val := expr.Path.Value
		if expr.Name != nil && expr.Name.String() != "" {
			val = expr.Name.String() + " " + val
		}
		f.Imports = append(f.Imports, val)
	case *ast.TypeSpec:
		tv := &TypeVisitor{
			Name: expr.Name.String(),
		}
		f.types = append(f.types, tv)
		return tv
	}
	return f
}

type TypeVisitor struct {
	Name  string
	Field []Field
}

func (t *TypeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch expr := node.(type) {
	case *ast.Field:
		var typ string
		switch ex := expr.Type.(type) {
		case *ast.Ident:
			typ = ex.String()
		case *ast.StarExpr:
			switch e := ex.X.(type) {
			case *ast.Ident:
				typ = "*" + e.String()
			case *ast.SelectorExpr:
				typ = "*" + e.X.(*ast.Ident).String() + "." + e.Sel.String()
			}
		case *ast.ArrayType:
			typ = "[]" + ex.Elt.(*ast.Ident).String()
		}
		for _, f := range expr.Names {
			t.Field = append(t.Field, Field{Name: f.String(), Typ: typ})
		}
	}
	return t
}

type Field struct {
	Name string
	Typ  string
}
