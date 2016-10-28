# langconv

langconv generate client side source code from golang source code.

## Install

```
go install github.com/nirasan/langconv
```

## How to use

```
langconv -f {GOLANG_FILENAME} -c {CONFIG_FILENAME} -o {OUTPUT_FILENAME}
```

## Example

### Golang source code

```go:sample.go
package main

const (
	CONST1 int32   = 1
	CONST2 string  = "hello hello hello"
	CONST3 float64 = 3.1412
)

type User struct {
	Username string
	Age      int
}
```

### Config file for C#

```toml:config.toml
ConstTemplate = '''
public static partial class Constant
{
{{ range . -}}
{{ "    " -}} public const {{ typeconv .Type }} {{ .Name }} = {{ .Value }};
{{ end -}}
}
'''

StructTemplate = '''
public class {{ .Name }}
{
{{ range .Fields -}}
{{ "    " -}} public {{ typeconv .Type }} {{ .Name }};
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
langconv -f sample.go -c config.toml -o sample.cs
```

### Output file

```cs
public static partial class Constant
{
    public const int CONST1 = 1;
    public const string CONST2 = "hello hello hello";
    public const double CONST3 = 3.1412;
}
public class User
{
    public string Username;
    public int Age;
}
```

## Other config example

### for Java

```toml
ConstTemplate = '''
public class Constant {
{{ range . -}}
{{ "    " -}} public static final {{ typeconv .Type }} {{ .Name }} = {{ .Value }};
{{ end -}}
}
'''

StructTemplate = '''
public class {{ .Name }} {
{{ range .Fields -}}
{{ "    " -}} public {{ typeconv .Type }} {{ .Name }};
{{ end -}}
}
'''

[Typemap]
int = "int"
int32 = "int"
int64 = "long"
float32 = "float"
float64 = "double"
string = "String"
bool = "boolean"
```
