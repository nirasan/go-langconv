package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestNewConstDecl(t *testing.T) {
	src := `
package main
const (
	CONST_INT int32 = 1
	CONST_STRING string = "hello"
)
const CONST_SINGLE bool = true
var VAR1 = 1
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	// ast.Print(fset, f)

	var cds []*ConstDecl
	var cd *ConstDecl

	cds = NewConstDecl(f.Decls[0])
	if len(cds) != 2 {
		t.Error("invalid length", len(cds))
	}
	cd = cds[0]
	if cd.Name != "CONST_INT" || cd.Type != "int32" || cd.Value != "1" {
		t.Error("invalid ConstDecl", cd.Name, cd.Type, cd.Value)
	}
	cd = cds[1]
	if cd.Name != "CONST_STRING" || cd.Type != "string" || cd.Value != `"hello"` {
		t.Error("invalid ConstDecl", cd.Name, cd.Type, cd.Value)
	}

	cds = NewConstDecl(f.Decls[1])
	if len(cds) != 1 {
		t.Error("invalid length", len(cds))
	}
	cd = cds[0]
	if cd.Name != "CONST_SINGLE" || cd.Type != "bool" || cd.Value != "true" {
		t.Error("invalid ConstDecl", cd.Name, cd.Type, cd.Value)
	}

	cds = NewConstDecl(f.Decls[2])
	if cds != nil {
		t.Error("cast error", cds)
	}
}

func TestNewStructDecl(t *testing.T) {
	src := `
package main
type User struct {
  Username string
  Age int32
  Email string
  EmailVerification bool
  Sights []float32
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	//ast.Print(fset, f)

	var sd *StructDecl
	sd = NewStructDecl(f.Decls[0])

	if sd == nil {
		t.Error("cast error")
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
