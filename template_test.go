package main

import "testing"

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

func TestRenderTemplate(t *testing.T) {
	cg := &ConstDeclGroup{
		IsEnum: false,
		ConstDeclList: []*ConstDecl{
			&ConstDecl{Name: "CONST1", Type: "int32", Value: "100"},
			&ConstDecl{Name: "CONST2", Type: "float32", Value: "3.14"},
			&ConstDecl{Name: "CONST3", Type: "string", Value: `"hello world"`},
		},
	}
	eg := &ConstDeclGroup{
		Name:   "MyEnum",
		IsEnum: true,
		ConstDeclList: []*ConstDecl{
			&ConstDecl{Name: "CONST1", Type: "int32", Value: "1"},
			&ConstDecl{Name: "CONST2", Type: "int32", Value: "2"},
			&ConstDecl{Name: "CONST3", Type: "int32", Value: "3"},
		},
	}
	s := &StructDecl{
		Name: "User",
		Fields: []StructDeclField{
			StructDeclField{Name: "Username", Type: "string", IsArray: false},
			StructDeclField{Name: "Age", Type: "int64", IsArray: false},
			StructDeclField{Name: "Sights", Type: "float32", IsArray: true},
		},
	}
	data := TemplateData{
		ClassName:              "MyClass",
		ConstDeclGroupList:     []*ConstDeclGroup{cg},
		EnumConstDeclGroupList: []*ConstDeclGroup{eg},
		StructDeclList:         []*StructDecl{s},
	}

	out := renderTemplate(data, defaultTemplate, defaultTypemap)

	expect := `public static class MyClass
{
    public const int CONST1 = 100;
    public const float CONST2 = 3.14;
    public const string CONST3 = "hello world";
    public enum MyEnum
    {
        CONST1 = 1,
        CONST2 = 2,
        CONST3 = 3,
    }
    public class User
    {
        public string Username;
        public long Age;
        public float[] Sights;
    }
}
`

	if out != expect {
		t.Errorf("invalid output\n%s\n\n%s", expect, out)
	}
}
