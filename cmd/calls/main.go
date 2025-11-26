package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const targetImportPath = "github.com/goxjs/gl"

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.Mode(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing directory: %v\n", err)
		return
	}

	uniques := make(map[string]struct{})
	sortedUniques := make([]string, 0)
	for pkgName, pkg := range pkgs {
		fmt.Printf("Analyzing package: **%s**\n", pkgName)

		glImports := getTargetImports(pkg, targetImportPath)

		if len(glImports) == 0 {
			fmt.Printf("No import of \"%s\" found in package **%s**.\n", targetImportPath, pkgName)
			continue
		}

		ast.Inspect(pkg, func(node ast.Node) bool {
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkgIdent, ok := selExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			if _, isTargetImport := glImports[pkgIdent.Name]; isTargetImport {
				callKey := fmt.Sprintf("%s.%s", pkgIdent.Name, selExpr.Sel.Name)
				fmt.Printf("  **Found call:** %s at %s\n", callKey, fset.Position(callExpr.Lparen))
				if _, found := uniques[callKey]; !found {
					uniques[callKey] = struct{}{}
					sortedUniques = append(sortedUniques, callKey)
				}
			}

			return true
		})
	}

	sort.Strings(sortedUniques)
	fmt.Printf("%s unique calls : \n", targetImportPath)
	for _, callKey := range sortedUniques {
		fmt.Printf("%s\n", callKey)
	}
}

func getTargetImports(pkg *ast.Package, targetPath string) map[string]*ast.ImportSpec {
	imports := make(map[string]*ast.ImportSpec)

	for _, file := range pkg.Files {
		for _, spec := range file.Imports {
			path := strings.Trim(spec.Path.Value, `"`)

			if path == targetPath {
				localName := ""
				if spec.Name != nil {
					localName = spec.Name.Name
				} else {
					_, localName = filepath.Split(path)
				}

				if localName != "_" && localName != "." {
					imports[localName] = spec
				}
			}
		}
	}
	return imports
}
