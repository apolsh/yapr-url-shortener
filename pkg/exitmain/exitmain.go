package exitmain

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer анализатор
var Analyzer = &analysis.Analyzer{
	Name: "exitmain",
	Doc:  "check for using os.Exit in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		funcDecl, ok := getMainFuncDecl(file)
		if !ok {
			continue
		}

		ast.Inspect(funcDecl, func(astNode ast.Node) bool {
			call, ok := astNode.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			var packageIdent *ast.Ident
			functionIdent := selector.Sel
			packageIdent, ok = selector.X.(*ast.Ident)
			if !ok {
				return true
			}
			if packageIdent.Name == "os" && functionIdent.Name == "Exit" {
				pass.Reportf(call.Pos(), "direct call os.Exit from main function")
			}

			return true
		})

	}
	return nil, nil
}

func getMainFuncDecl(file *ast.File) (*ast.FuncDecl, bool) {
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == "main" {
				return funcDecl, true
			}
		}
	}
	return nil, false
}
