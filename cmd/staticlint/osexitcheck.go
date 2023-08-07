package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var osExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit call in main function within package main",
	Run:  osExitRun,
}

func osExitRun(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "main" {
				return true
			}

			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				x, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				obj, ok := pass.TypesInfo.Uses[x].(*types.PkgName)
				if !ok || obj.Imported().Path() != "os" || sel.Sel.Name != "Exit" {
					return true
				}

				pass.Reportf(callExpr.Pos(), "direct call to os.Exit in main function of main package")
				return true
			})

			return true
		})
	}
	return nil, nil

}
