package main

import (
	"bytes"
	"html/template"
	"log"
)

type TemplateData struct {
	ClassName              string
	ConstDeclGroupList     []*ConstDeclGroup
	EnumConstDeclGroupList []*ConstDeclGroup
	StructDeclList         []*StructDecl
}

func renderTemplate(data TemplateData, tmpltext string, typemap map[string]string) string {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"typeconv": func(t string) string {
			v, ok := typemap[t]
			if ok {
				return v
			} else {
				return t
			}
		},
	}).Parse(tmpltext))
	var buf bytes.Buffer
	e := tmpl.Execute(&buf, data)
	if e != nil {
		log.Fatal(e)
	}
	return buf.String()
}
