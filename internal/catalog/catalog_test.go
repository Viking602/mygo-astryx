package catalog

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// Gallery coverage is a product contract: a new public component must not be
// silently absent from its index. App tests then verify each index has a demo.
func TestCatalogTracksPublicComponents(t *testing.T) {
	files, err := filepath.Glob("../../*.go")
	if err != nil {
		t.Fatal(err)
	}
	methods := map[string]bool{}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || !fn.Name.IsExported() {
				continue
			}
			pointer, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			receiver, ok := pointer.X.(*ast.Ident)
			if !ok || receiver.Name != "UI" {
				continue
			}
			if fn.Name.Name == "Tokens" || fn.Name.Name == "Key" {
				continue
			}
			methods[fn.Name.Name] = true
		}
	}
	for _, component := range Components {
		if !methods[component.Name] {
			t.Errorf("catalogued component has no public method: %s", component.Name)
		}
		delete(methods, component.Name)
	}
	for name := range methods {
		t.Errorf("public component missing from Gallery catalog: %s", name)
	}
}
