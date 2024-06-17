package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSExitAnalyzer - analyzer of unwanted usage of os.Exit call
var OSExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for usage of os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			// Не обходим не main функции
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}
			}
			if call, ok := node.(*ast.CallExpr); ok {
				var id *ast.Ident
				switch fun := call.Fun.(type) {
				case *ast.Ident:
					id = fun
				case *ast.SelectorExpr:
					id = fun.Sel
				}
				if id != nil && !pass.TypesInfo.Types[id].IsType() && id.Name == "Exit" {
					pass.Report(analysis.Diagnostic{
						Pos:     call.Lparen,
						Message: fmt.Sprintf("prohibited call of %s(...)", id.Name),
					})
				}
			}
			return true
		})
	}
	return nil, nil
}
