package interputil_test

import (
	"testing"

	"github.com/axw/interputil"
)

func TestExplicitNewline(t *testing.T) {
	var b interputil.Buffer
	n, err := b.WriteString("what\n")
	if err != nil {
		t.Fatalf("failed to write partial raw string: %v", err)
	}
	if n != 5 {
		t.Fatalf("expected 5, got %d", n)
	}
}

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

func TestWriteMultipleLines(t *testing.T) {
	var b interputil.Buffer
	n, err := b.WriteString("what\never")
	if err != interputil.ErrMultipleLines {
		t.Fatalf("expected %v, got %v", interputil.ErrMultipleLines, err)
	}
	if n != len("what\n") {
		t.Fatalf("expected %v, got %v", len("what\n"), n)
	}
	if !b.Ready() {
		t.Errorf("expected buffer to be ready")
	}
	if s := b.String(); s != "what\n" {
		t.Errorf("expected what\n, got %q", s)
	}
}

// TODO
func TestComments(t *testing.T) {
	var b interputil.Buffer
	if _, err := b.WriteString("// abc\nmeep"); err != nil {
		t.Fatal(err)
	}
	t.Log(b.String())
	t.Logf("%v", b.Tokens())
}
