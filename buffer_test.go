package interputil_test

import (
	"testing"

	"llvm.org/llvm/tools/llgo/interp/interputil"
)

func TestBackQuotes(t *testing.T) {
	var b interputil.Buffer
	if _, err := b.WriteString("`abc"); err != nil {
		t.Fatalf("failed to write partial raw string: %v", err)
	}
	if b.Ready() {
		t.Errorf("expected buffer not to be ready")
	}
	if _, err := b.WriteString("def`"); err != nil {
		t.Fatalf("failed to complete writing raw string: %v", err)
	}
	if !b.Ready() {
		t.Errorf("expected buffer to be ready")
	}
	if s := b.String(); s != "`abcdef`" {
		t.Errorf("expected `abcdef`, got %q", s)
	}
}
