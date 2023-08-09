package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"golang.org/x/tools/go/analysis"
)

var osExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit call in main function within package main",
	Run:  osExitRun,
}

func osExitRun(pass *analysis.Pass) (interface{}, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "../shortener/main.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "os" && fun.Sel.Name == "Exit" {
					log.Fatal("os.Exit is not allowed in main function")
				}
			}
		}
		return true
	})
	return nil, nil
}
