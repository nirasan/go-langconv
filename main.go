package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"bufio"
)

var classFlag = flag.String("class", "LangConv", "output file class")
var dirFlag = flag.String("d", "", "source dir")
var configFlag = flag.String("c", "", "config file")
var outputFlag = flag.String("o", "", "output file")

var constDeclGroupList []*ConstDeclGroup
var structDeclList []*StructDecl

func main() {
	flag.Parse()

	// load config
	config := loadConfig(*configFlag)

	// parse go files
	constDeclGroupList = []*ConstDeclGroup{}
	structDeclList = []*StructDecl{}
	filepath.Walk(*dirFlag, walker)

	// create template data
	data := TemplateData{
		ClassName:              *classFlag,
		ConstDeclGroupList:     []*ConstDeclGroup{},
		EnumConstDeclGroupList: []*ConstDeclGroup{},
		StructDeclList:         []*StructDecl{},
	}
	for i := len(constDeclGroupList) - 1; i >= 0; i-- {
		g := constDeclGroupList[i]
		if g.IsEnum {
			data.EnumConstDeclGroupList = append(data.EnumConstDeclGroupList, g)
		} else {
			data.ConstDeclGroupList = append(data.ConstDeclGroupList, g)
		}
	}
	data.StructDeclList = structDeclList

	// write output file
	fp := newFile(*outputFlag)
	defer fp.Close()
	out := renderTemplate(data, config.Template, config.Typemap)
	w := bufio.NewWriter(fp)
	if _, e := w.WriteString(out); e != nil {
		log.Fatal(e)
	}
	w.Flush()
}

func walker(path string, info os.FileInfo, err error) error {
	if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
		return nil
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, decl := range f.Decls {
		tdecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch tdecl.Tok {
		case token.CONST:
			if g := NewConstDeclGroup(decl); g != nil {
				constDeclGroupList = append(constDeclGroupList, g)
			}
		case token.TYPE:
			if sd := NewStructDecl(decl); sd != nil {
				structDeclList = append(structDeclList, sd)
			}
		}
	}
	return nil
}

func newFile(fn string) *os.File {
	fp, err := os.OpenFile(fn, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
