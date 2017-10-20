package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"
)

type generator struct {
	buf    bytes.Buffer
	indent string

	filename                  string // may be empty
	srcPackage, srcInterfaces string // may be empty

	packageMap map[string]string // map from import path to package name
}

func main() {
	g := generator{
		filename: "sample.go",
	}
	names := parseFuncName("sample.go")

	g.Generate(names)
	fltrName := strings.Split(g.filename, ".")[0]
	ioutil.WriteFile(fltrName+"Test.go", g.Output(), 0644)
}

func parseFuncName(fileName string) []string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var funcNames []string
	fmt.Println("Functions:")
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcNames = append(funcNames, fn.Name.Name)
		fmt.Println(fn.Name.Name)
	}

	return funcNames
}

func (g *generator) p(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, g.indent+format+"\n", args...)
}

func (g *generator) Generate(funcNames []string) {
	for _, k := range funcNames {
		g.p("func Test%v(t testing.T){", k)
		g.p("}")
	}
}

func (g *generator) Output() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Fatalf("Failed to format generated source code: %s\n%s", err, g.buf.String())
	}
	return src
}
