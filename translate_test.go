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

func TestTranslateConst(t *testing.T) {
	consts := []*ConstDecl{
		&ConstDecl{Name: "CONST1", Type: "int32", Value: "100"},
		&ConstDecl{Name: "CONST2", Type: "float32", Value: "3.14"},
		&ConstDecl{Name: "CONST3", Type: "string", Value: `"hello world"`},
	}
	tmpl := `
public static partial class Constant
{
{{ range . -}}
{{ "    " -}} public const {{ typeconv .Type }} {{ .Name }} = {{ .Value }};
{{ end -}}
}
`
	out := TranslateConst(consts, tmpl, typemap)

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
