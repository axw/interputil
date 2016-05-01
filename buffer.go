package interputil

import (
	"bytes"
	"go/scanner"
	"go/token"
	"strings"
)

// Buffer encapsulates a string buffer, keep tracking of whether the buffer
// contains a complete declaration, statement, or expression.
type Buffer struct {
	buf   bytes.Buffer
	first token.Token

	// parens, bracks, and braces count the number of open parentheses,
	// brackets and braces, so we know when they are balanced.
	parens, bracks, braces int

	// backquote records whether or not a back quoted string has been
	// started, but not completed.
	backquote bool
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
//
// WriteString will also track the number of open parentheses, brackets,
// braces, and back quotes. Whenever these are all balanced, the Ready
// method will return true.
//
// WriteString currently assumes that you will not write more than one
// declaration, statement, or expression at a time.
func (b *Buffer) WriteString(s string) (int, error) {
	var scanner scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", -1, len(s))
	scanner.Init(file, []byte(s), nil, 0)
	scan := scanner.Scan

	for _, tok, _ := scan(); tok != token.EOF; _, tok, _ = scan() {
		if b.first == 0 {
			b.first = tok
		}
		switch tok {
		case token.LPAREN:
			b.parens++
		case token.RPAREN:
			b.parens--
		case token.LBRACE:
			b.braces++
		case token.RBRACE:
			b.braces--
		case token.LBRACK:
			b.bracks++
		case token.RBRACK:
			b.bracks--
		case token.STRING:
			if strings.HasPrefix(s, "`") {
				b.backquote = !b.backquote
			}
			if len(s) > 1 && strings.HasSuffix(s, "`") {
				b.backquote = !b.backquote
			}
		}
	}

	return b.buf.WriteString(s)
}

// Ready returns true iff the buffer contains a balanced set of parentheses,
// brackets, braces, and back quotes.
func (b *Buffer) Ready() bool {
	return b.parens <= 0 && b.bracks <= 0 && b.braces <= 0 && !b.backquote
}

// First returns the first token in the buffer. If nothing has been added to
// the buffer, First will return token.ILLEGAL.
func (b *Buffer) First() token.Token {
	return b.first
}

// String returns the contents of the buffer as a string. If the Buffer is a
// nil pointer, it returns "<nil>".
func (b *Buffer) String() string {
	return b.buf.String()
}
