package interputil

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseImports(buf *Buffer) ([]*ast.ImportSpec, error) {
	if buf.First() != token.IMPORT {
		panic("buffer does not contain an import spec")
	}
	file, err := parseTopLevel(buf)
	if err != nil {
		return nil, err
	}
	return file.Imports, nil
}

func ParseFuncDecl(buf *Buffer) (*ast.FuncDecl, error) {
	if buf.First() != token.FUNC {
		panic("buffer does not contain a type spec")
	}
	file, err := parseTopLevel(buf)
	if err != nil {
		return nil, err
	}
	return file.Decls[0].(*ast.FuncDecl), nil
}

func ParseTypeSpec(buf *Buffer) (*ast.TypeSpec, error) {
	if buf.First() != token.TYPE {
		panic("buffer does not contain a type spec")
	}
	file, err := parseTopLevel(buf)
	if err != nil {
		return nil, err
	}
	return file.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec), nil
}

func ParseValueSpec(buf *Buffer) (*ast.ValueSpec, error) {
	if t := buf.First(); t != token.VAR && t != token.CONST {
		panic("buffer does not contain a value spec")
	}
	file, err := parseTopLevel(buf)
	if err != nil {
		return nil, err
	}
	return file.Decls[0].(*ast.GenDecl).Specs[0].(*ast.ValueSpec), nil
}

func ParseStmt(buf *Buffer) (ast.Stmt, error) {
	src := []byte("package p\nfunc f(){" + buf.String() + "}")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "<input>", src, parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return file.Decls[0].(*ast.FuncDecl).Body.List[0], nil
}

// parseTopLevel parses a single top-level declaration, returning an *ast.File
// containing it.
func parseTopLevel(buf *Buffer) (*ast.File, error) {
	src := []byte("package p\n" + buf.String())
	fset := token.NewFileSet()
	return parser.ParseFile(fset, "<input>", src, parser.DeclarationErrors|parser.ParseComments)
}
