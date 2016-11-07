package main

import (
	"bytes"
	"log"
	"text/template"
)

func TranslateConst(c []*ConstDecl, tmpltext string, typemap map[string]string) string {
	return translateBase(c, tmpltext, typemap)
}

func TranslateStruct(s *StructDecl, tmpltext string, typemap map[string]string) string {
	return translateBase(s, tmpltext, typemap)
}

func translateBase(data interface{}, tmpltext string, typemap map[string]string) string {
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
