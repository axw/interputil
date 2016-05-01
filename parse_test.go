package interputil_test

import (
	"go/ast"
	"testing"

	"github.com/axw/interputil"
)

func TestParseImports(t *testing.T) {
	testParseImports(t, `import "abc"`, importSpec{"", `"abc"`})
	testParseImports(t, `import woop "abc"`, importSpec{"woop", `"abc"`})
	testParseImports(t, `import . "abc"`, importSpec{".", `"abc"`})
	testParseImports(t, `import _ "abc"`, importSpec{"_", `"abc"`})
	testParseImports(t, `import ()`)
	testParseImports(t, `import (
  "abc"
  . "def"
)`, importSpec{"", `"abc"`}, importSpec{".", `"def"`})
}

func testParseImports(t *testing.T, src string, expect ...importSpec) {
	t.Logf("testing ParseImports with %q", src)
	b := newBuffer(src)
	imports, err := interputil.ParseImports(b)
	if err != nil {
		t.Errorf("failed to parse %q: %v", src, err)
		return
	}
	if len(imports) != len(expect) {
		t.Errorf("expected %d imports, got %d", len(expect), len(imports))
		return
	}
	for i, spec := range imports {
		var actual importSpec
		if spec.Name != nil {
			actual.name = spec.Name.Name
		}
		actual.path = spec.Path.Value
		if actual != expect[i] {
			t.Errorf("import spec %d mismatch: expected %v, got %v", expect[i], actual)
			return
		}
	}
}

type importSpec struct {
	name string
	path string
}

func TestParseFuncDecl(t *testing.T) {
	testParseFuncDecl(t, "func F1() {\n}", "F1")
	testParseFuncDecl(t, "func f2(arg int) {}", "f2")
	testParseFuncDecl(t, "func F(a1 int, a2 byte) (r1 int, r2 string) {}", "F")
}

func TestParseFuncDeclRecv(t *testing.T) {
	f := testParseFuncDecl(t, "func (t T) F() {}", "F")
	if f.Recv == nil || f.Recv.NumFields() != 1 {
		t.Fatal("expected non-nil receiver")
	}
	// Methods are rejected by llgoi on the basis that they have a
	// receiver, so we've checked all we need to check here.
}

func testParseFuncDecl(t *testing.T, src, name string) *ast.FuncDecl {
	t.Logf("testing ParseFuncDecl with %q", src)
	b := newBuffer(src)
	funcDecl, err := interputil.ParseFuncDecl(b)
	if err != nil {
		t.Errorf("failed to parse %q: %v", src, err)
		return nil
	}
	if funcDecl.Name.Name != name {
		t.Errorf("expected function name %q, got %q", name, funcDecl.Name.Name)
		return nil
	}
	return funcDecl
}

func TestParseTypeSpec(t *testing.T) {
	testParseTypeSpec(t, "type X int")
	testParseTypeSpec(t, "type X struct{Y int\nZ string\n}")
}

func testParseTypeSpec(t *testing.T, src string) *ast.TypeSpec {
	t.Logf("testing ParseTypeSpec with %q", src)
	b := newBuffer(src)
	spec, err := interputil.ParseTypeSpec(b)
	if err != nil {
		t.Errorf("failed to parse %q: %v", src, err)
		return nil
	}
	return spec
}

func TestParseValueSpec(t *testing.T) {
	spec := testParseValueSpec(t, "var X int")
	if n := len(spec.Names); n != 1 {
		t.Errorf("expected 1 name, got %d", n)
		return
	}
	if name := spec.Names[0].Name; name != "X" {
		t.Errorf("expected name of %q, got %q", "X", name)
	}
}

func testParseValueSpec(t *testing.T, src string) *ast.ValueSpec {
	t.Logf("testing ParseValueSpec with %q", src)
	b := newBuffer(src)
	spec, err := interputil.ParseValueSpec(b)
	if err != nil {
		t.Errorf("failed to parse %q: %v", src, err)
		return nil
	}
	return spec
}

func TestParseStmt(t *testing.T) {
	testParseStmt(t, "x := 123")
	testParseStmt(t, "x++")
	testParseStmt(t, "x := func(){\n}")
	testParseStmt(t, "1+2")
}

func testParseStmt(t *testing.T, src string) ast.Stmt {
	t.Logf("testing ParseStmt with %q", src)
	b := newBuffer(src)
	stmt, err := interputil.ParseStmt(b)
	if err != nil {
		t.Errorf("failed to parse %q: %v", src, err)
		return nil
	}
	return stmt
}

func newBuffer(str string) *interputil.Buffer {
	var b interputil.Buffer
	if _, err := b.WriteString(str); err != nil {
		panic(err)
	}
	return &b
}
