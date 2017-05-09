package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"bufio"

	"github.com/naoina/toml"
)

var dirFlag = flag.String("d", "", "source dir")
var configFlag = flag.String("c", "", "config file")
var outputFlag = flag.String("o", "", "output file")

var constDeclList [][]*ConstDecl
var structDeclList []*StructDecl

type Config struct {
	ConstTemplate  string
	StructTemplate string
	Typemap        map[string]string
}

func main() {
	flag.Parse()

	// load config
	config := loadConfig()

	// parse go files
	constDeclList = [][]*ConstDecl{}
	structDeclList = []*StructDecl{}
	filepath.Walk(*dirFlag, walker)

	// write output file
	fp := newFile(*outputFlag)
	defer fp.Close()
	w := bufio.NewWriter(fp)

	for i := len(constDeclList) - 1; i >= 0; i-- {
		cd := constDeclList[i]
		s := TranslateConst(cd, config.ConstTemplate, config.Typemap)
		if _, e := w.WriteString(s); e != nil {
			log.Fatal(e)
		}
	}

	for i := len(structDeclList) - 1; i >= 0; i-- {
		sd := structDeclList[i]
		s := TranslateStruct(sd, config.StructTemplate, config.Typemap)
		if _, e := w.WriteString(s); e != nil {
			log.Fatal(e)
		}
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
			if cd := NewConstDecl(decl); cd != nil {
				constDeclList = append(constDeclList, cd)
			}
		case token.TYPE:
			if sd := NewStructDecl(decl); sd != nil {
				structDeclList = append(structDeclList, sd)
			}
		}
	}
	return nil
}

func loadConfig() Config {
	confData, err := os.Open(*configFlag)
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(confData)
	if err != nil {
		panic(err)
	}
	var config Config
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	return config
}

func newFile(fn string) *os.File {
	fp, err := os.OpenFile(fn, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
