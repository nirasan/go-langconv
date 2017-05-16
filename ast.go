package main

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

const CommentPrefix = "+langconv"

type ConstDeclGroup struct {
	Name          string
	IsEnum        bool
	ConstDeclList []*ConstDecl
}

type ConstDecl struct {
	Name  string
	Type  string
	Value string
}

func NewConstDeclGroup(decl ast.Decl) *ConstDeclGroup {
	// validate const decl
	tdecl, ok := decl.(*ast.GenDecl)
	if !ok || tdecl.Tok != token.CONST || tdecl.Doc == nil {
		return nil
	}
	// validate comment
	comment := tdecl.Doc.Text()
	if strings.Index(comment, CommentPrefix) == -1 {
		return nil
	}
	// create ConstDeclGroup
	g := &ConstDeclGroup{
		ConstDeclList: []*ConstDecl{},
	}
	// parse enum option
	r := regexp.MustCompile(`\senum:(\S+)\s`)
	matched := r.FindAllStringSubmatch(comment, -1)
	if len(matched) > 0 {
		g.IsEnum = true
		g.Name = matched[0][1]
	}
	// create ConstDecl
	for _, s := range tdecl.Specs {
		ts, ok := s.(*ast.ValueSpec)
		if !ok {
			continue
		}
		// get name and new ConstDecl
		c := &ConstDecl{Name: ts.Names[0].Name}
		// get type if defined
		if ts.Type != nil {
			if t, ok := ts.Type.(*ast.Ident); ok {
				c.Type = t.Name
			}
		}
		// get value
		if ts.Values != nil {
			switch v := ts.Values[0].(type) {
			case *ast.BasicLit: // like int and string
				c.Value = v.Value
			case *ast.Ident: // like bool
				c.Value = v.Name
			}
		}
		g.ConstDeclList = append(g.ConstDeclList, c)
	}
	return g
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
	// validate TypeDecl
	tdecl, ok := decl.(*ast.GenDecl)
	if !ok || tdecl.Tok != token.TYPE {
		return nil
	}
	// validate comment
	comment := tdecl.Doc.Text()
	if strings.Index(comment, CommentPrefix) == -1 {
		return nil
	}
	// validate type
	s, ok := tdecl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	// validate struct
	t, ok := s.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	// get name and new StructDecl
	sd := &StructDecl{
		Name:   s.Name.Name,
		Fields: []StructDeclField{},
	}
	// get fields
	for _, f := range t.Fields.List {
		if f.Names == nil {
			continue
		}
		// get field name
		sdf := StructDeclField{
			Name: f.Names[0].Name,
		}
		// get field type
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
