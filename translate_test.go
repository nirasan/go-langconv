package main

import (
	"testing"
)

var typemap = map[string]string{
	"int":     "int",
	"int32":   "int",
	"int64":   "long",
	"uint":    "uint",
	"uint32":  "uint",
	"uint64":  "ulong",
	"float32": "float",
	"float64": "double",
	"string":  "string",
	"bool":    "bool",
}

func TestTranslateConstGroup(t *testing.T) {
	consts := []*ConstDecl{
		&ConstDecl{Name: "CONST1", Type: "int32", Value: "100"},
		&ConstDecl{Name: "CONST2", Type: "float32", Value: "3.14"},
		&ConstDecl{Name: "CONST3", Type: "string", Value: `"hello world"`},
	}
	g := &ConstDeclGroup{
		ConstDeclList: consts,
	}
	tmpl := `
public static partial class Constant
{
{{ range .ConstDeclList -}}
{{ "    " -}} public const {{ typeconv .Type }} {{ .Name }} = {{ .Value }};
{{ end -}}
}
`
	out := TranslateConstGroup(g, tmpl, typemap)

	expect := `
public static partial class Constant
{
    public const int CONST1 = 100;
    public const float CONST2 = 3.14;
    public const string CONST3 = "hello world";
}
`

	if out != expect {
		t.Error("invalid output", out)
	}
}

func TestTranslateConstGroup2(t *testing.T) {
	consts := []*ConstDecl{
		&ConstDecl{Name: "CONST1", Type: "int32", Value: "1"},
		&ConstDecl{Name: "CONST2", Type: "int32", Value: "2"},
		&ConstDecl{Name: "CONST3", Type: "int32", Value: "3"},
	}
	g := &ConstDeclGroup{
		Name:          "MyEnum",
		IsEnum:        true,
		ConstDeclList: consts,
	}
	tmpl := `
public static partial class Constant
{
    public enum {{ .Name }}
    {
{{ range .ConstDeclList -}}
{{ "        " -}} {{ .Name }} = {{ .Value }},
{{ end -}}
{{ "    " -}} }
}
`
	out := TranslateConstGroup(g, tmpl, typemap)

	expect := `
public static partial class Constant
{
    public enum MyEnum
    {
        CONST1 = 1,
        CONST2 = 2,
        CONST3 = 3,
    }
}
`

	if out != expect {
		t.Errorf("invalid output\n\n%s\n\n%s\n", expect, out)
	}
}

func TestTranslateStruct(t *testing.T) {
	s := &StructDecl{
		Name: "User",
		Fields: []StructDeclField{
			StructDeclField{Name: "Username", Type: "string", IsArray: false},
			StructDeclField{Name: "Age", Type: "int64", IsArray: false},
			StructDeclField{Name: "Sights", Type: "float32", IsArray: true},
		},
	}

	tmpl := `
public class {{ .Name }}
{
{{ range .Fields -}}
{{ "    " -}} public {{ typeconv .Type }} {{- if .IsArray -}} [] {{- end }} {{ .Name }};
{{ end -}}
}
`

	out := TranslateStruct(s, tmpl, typemap)

	expect := `
public class User
{
    public string Username;
    public long Age;
    public float[] Sights;
}
`
	if out != expect {
		t.Error("invalid output", out)
	}
}
