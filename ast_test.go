package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNewConstDeclGroup(t *testing.T) {
	src := `
package main
// +langconv
const (
	CONST_INT int32 = 1
	CONST_STRING string = "hello"
)
// +langconv enum:MyEnum
const CONST_SINGLE bool = true
// +langconv
var VAR1 = 1
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err.Error())
	}

	ast.Print(fset, f)

	var g *ConstDeclGroup
	var cd *ConstDecl

	g = NewConstDeclGroup(f.Decls[0])
	if len(g.ConstDeclList) != 2 || g.IsEnum {
		t.Fatalf("invalid struct: %v", g)
	}
	cd = g.ConstDeclList[0]
	if cd.Name != "CONST_INT" || cd.Type != "int32" || cd.Value != "1" {
		t.Error("invalid ConstDecl", cd.Name, cd.Type, cd.Value)
	}
	cd = g.ConstDeclList[1]
	if cd.Name != "CONST_STRING" || cd.Type != "string" || cd.Value != `"hello"` {
		t.Error("invalid ConstDecl", cd.Name, cd.Type, cd.Value)
	}

	g = NewConstDeclGroup(f.Decls[1])
	if len(g.ConstDeclList) != 1 || !g.IsEnum || g.Name != "MyEnum" {
		t.Fatalf("invalid struct: %+v", g)
	}
	cd = g.ConstDeclList[0]
	if cd.Name != "CONST_SINGLE" || cd.Type != "bool" || cd.Value != "true" {
		t.Error("invalid ConstDecl", cd.Name, cd.Type, cd.Value)
	}

	g = NewConstDeclGroup(f.Decls[2])
	if g != nil {
		t.Error("cast error", g)
	}
}

func TestNewStructDecl(t *testing.T) {
	src := `
package main
// +langconv
type User struct {
  Username string
  Age int32
  Email string
  EmailVerification bool
  Sights []float32
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err.Error())
	}

	ast.Print(fset, f)

	var sd *StructDecl
	sd = NewStructDecl(f.Decls[0])

	if sd == nil {
		t.Fatal("parse error")
	}
	if sd.Name != "User" {
		t.Error("invalid name", sd.Name)
	}
	samples := []struct {
		Name    string
		Type    string
		IsArray bool
	}{
		{"Username", "string", false},
		{"Age", "int32", false},
		{"Email", "string", false},
		{"EmailVerification", "bool", false},
		{"Sights", "float32", true},
	}
	for i, sample := range samples {
		if sd.Fields[i].Name != sample.Name || sd.Fields[i].Type != sample.Type {
			t.Error("invalid field", sd.Fields[i].Name, sd.Fields[i].Type)
		}
	}
}
