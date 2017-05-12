package main

import (
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Template string
	Typemap  map[string]string
}

var defaultTemplate = `public static class {{ .ClassName }}
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
`

var defaultTypemap = map[string]string{
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

func loadConfig(filename string) Config {
	config := Config{}
	if filename != "" {
		data, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		buf, err := ioutil.ReadAll(data)
		if err != nil {
			panic(err)
		}
		if err := toml.Unmarshal(buf, &config); err != nil {
			panic(err)
		}
	}
	if config.Template == "" {
		config.Template = defaultTemplate
	}
	if len(config.Typemap) == 0 {
		config.Typemap = defaultTypemap
	}
	return config
}
