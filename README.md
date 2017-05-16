# langconv

langconv convert golang source code into client side source code.

## Install

```
go install github.com/nirasan/go-langconv
```

## How to use

```
go-langconv -d {DIRNAME} -c {CONFIG_FILENAME} -o {OUTPUT_FILENAME} -class {OUTPUT_CLASSNAME}
```

## Example

### Golang source code

```go
package main

// +langconv
const (
	CONST1 int32   = 1
	CONST2 string  = "hello hello hello"
	CONST3 float64 = 3.1412
)

// +langconv enum:ItemCategory
const (
    Category1 = 1
    Category2 = 2
    Category3 = 3
)

// +langconv
type User struct {
	Username string
	Age      int
}
```

### Config file for C Sharp

```toml
Template = '''
public static class {{ .ClassName }}
{
{{ range .ConstDeclGroupList -}}
{{ range .ConstDeclList -}}
{{ "    " -}} public const {{ typeconv .Type }} {{ .Name }} = {{ .Value }};
{{ end -}}
{{ end -}}

{{ range .EnumConstDeclGroupList -}}
{{ "    public enum " -}} {{ .Name }}
{{ "    " -}} {
{{ range .ConstDeclList -}}
{{ "        " -}} {{ .Name }} = {{ .Value }},
{{ end -}}
{{ "    " -}} }
{{ end -}}

{{ range .StructDeclList -}}
{{ "    public class " -}} {{ .Name }}
{{ "    " -}} {
{{ range .Fields -}}
{{ "        " -}} public {{ typeconv .Type }} {{- if .IsArray -}} [] {{- end }} {{ .Name }};
{{ end -}}
{{ "    " -}} }
{{ end -}}
}
'''

[Typemap]
int = "int"
int32 = "int"
int64 = "long"
uint = "uint"
uint32 = "uint"
uint64 = "ulong"
float32 = "float"
float64 = "double"
string = "string"
bool = "bool"
```

### Run command

```
go-langconv -d dir -c config.toml -o sample.cs -class MyClass
```

### Output file

```cs
public static class MyClass
{
    public const int CONST1 = 1;
    public const string CONST2 = "hello hello hello";
    public const double CONST3 = 3.1412;
    public enum ItemCategory
    {
        Category1 = 1,
        Category2 = 2,
        Category3 = 3,
    }
    public class User
    {
        public string Username;
        public long Age;
    }
}
```
