package main

import (
	"go/ast"
	"go/token"
)

type ConstDecl struct {
	Name  string
	Type  string
	Value string
}

func NewConstDecl(decl ast.Decl) []*ConstDecl {
	tdecl, ok := decl.(*ast.GenDecl)
	if !ok || tdecl.Tok != token.CONST {
		return nil
	}
	consts := []*ConstDecl{}
	for _, s := range tdecl.Specs {
		ts, ok := s.(*ast.ValueSpec)
		if !ok {
			continue
		}
		c := &ConstDecl{Name: ts.Names[0].Name}
		if ts.Type != nil {
			if t, ok := ts.Type.(*ast.Ident); ok {
				c.Type = t.Name
			}
		}
		if ts.Values != nil {
			switch v := ts.Values[0].(type) {
			case *ast.BasicLit:
				c.Value = v.Value
			case *ast.Ident:
				c.Value = v.Name
			}
		}
		consts = append(consts, c)
	}
	return consts
}

type StructDecl struct {
	Name   string
	Fields []StructDeclField
}

type StructDeclField struct {
	Name    string
	Type    string
	IsArray bool
}

func NewStructDecl(decl ast.Decl) *StructDecl {
	tdecl, ok := decl.(*ast.GenDecl)
	if !ok || tdecl.Tok != token.TYPE {
		return nil
	}
	s, ok := tdecl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	t, ok := s.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	sd := &StructDecl{
		Name:   s.Name.Name,
		Fields: []StructDeclField{},
	}
	for _, f := range t.Fields.List {
		sdf := StructDeclField{
			Name: f.Names[0].Name,
		}
		switch ft := f.Type.(type) {
		case *ast.Ident:
			sdf.Type = ft.Name
		case *ast.ArrayType:
			sdf.IsArray = true
			sdf.Type = ft.Elt.(*ast.Ident).Name
		}
		sd.Fields = append(sd.Fields, sdf)
	}
	return sd
}
