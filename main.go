package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"

	"bufio"

	"github.com/naoina/toml"
)

var fFlag = flag.String("f", "", "source file")
var cFlag = flag.String("c", "", "config file")
var oFlag = flag.String("o", "", "output file")

type Config struct {
	ConstTemplate  string
	StructTemplate string
	Typemap        map[string]string
}

func main() {
	flag.Parse()
	config := loadConfig()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *fFlag, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}
	fp := newFile(*oFlag)
	defer fp.Close()
	w := bufio.NewWriter(fp)
	for _, decl := range f.Decls {
		tdecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch tdecl.Tok {
		case token.CONST:
			cd := NewConstDecl(decl)
			if cd == nil {
				continue
			}
			s := TranslateConst(cd, config.ConstTemplate, config.Typemap)
			if _, e := w.WriteString(s); e != nil {
				log.Fatal(e)
			}
		case token.TYPE:
			sd := NewStructDecl(decl)
			if sd == nil {
				continue
			}
			s := TranslateStruct(sd, config.StructTemplate, config.Typemap)
			if _, e := w.WriteString(s); e != nil {
				log.Fatal(e)
			}
		}
	}
	w.Flush()
}

func loadConfig() Config {
	confData, err := os.Open(*cFlag)
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
