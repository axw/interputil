package interputil

import (
	"go/scanner"
	"go/token"
	"strings"

	"github.com/juju/errors"
)

// ErrMultipleLines is returned from Buffer.WriteString if there is additional
// text after the end of the line. The caller should add this to the buffer
// again after consuming the complete, buffered line.
var ErrMultipleLines = errors.New("additional lines ignored")

// Buffer encapsulates a string buffer, keep tracking of whether the buffer
// contains a complete declaration, statement, or expression.
type Buffer struct {
	// complete contains a complete "line", up to and including the
	// newline.
	complete string

	// incomplete contains anything written beyond the end of the
	// complete line, or a fragment of an incomplete line.
	incomplete string

	// tokens contains the tokens that make up the complete line.
	tokens []token.Token
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
//
// WriteString will also track the number of open parentheses, brackets,
// braces, and back quotes. Whenever these are all balanced, the Ready
// method will return true.
//
// Comments at the beginning of a line will be ignored, to provide the
// invariant that the first token's offset aligns with the beginning of
// the buffer's contents.
func (b *Buffer) WriteString(s string) (int, error) {
	s = b.incomplete + s
	b.incomplete = s
	b.complete = ""
	b.tokens = nil
	var skipped int

	var scanner scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", -1, len(s))
	scanner.Init(file, []byte(s), nil, 0)
	scan := scanner.Scan

	var parens, bracks, braces int
	var backquote bool
	var tokens []token.Token

	pos, tok, lit := scan()
	beginPos := pos

loop:
	for ; tok != token.EOF; pos, tok, lit = scan() {
		if tok == token.ILLEGAL {
			return -1, errors.Errorf("illegal token: %v", lit)
		}

		switch tok {
		case token.SEMICOLON:
			if !backquote && parens+bracks+braces == 0 {
				beginOffset := 0
				if len(tokens) > 0 {
					beginOffset = fset.Position(beginPos).Offset
					skipped = beginOffset
				}
				endOffset := fset.Position(pos).Offset
				if n := endOffset + len(lit); n <= len(s) {
					endOffset = n
				}
				b.complete = s[beginOffset:endOffset]
				b.incomplete = ""
				b.tokens = tokens
				break loop
			}
		case token.LPAREN:
			parens++
		case token.RPAREN:
			parens--
		case token.LBRACE:
			braces++
		case token.RBRACE:
			braces--
		case token.LBRACK:
			bracks++
		case token.RBRACK:
			bracks--
		case token.STRING:
			if strings.HasPrefix(s, "`") {
				backquote = !backquote
			}
			if len(s) > 1 && strings.HasSuffix(s, "`") {
				backquote = !backquote
			}
		}
		tokens = append(tokens, tok)
	}

	if n := len(b.complete) + len(b.incomplete) + skipped; n != len(s) {
		return n, ErrMultipleLines
	}
	return len(s), nil
}

// Len returns the number of bytes in the buffer. Note that the length returned
// may exceed the len(b.String()), as the Len() includes the incomplete data.
func (b *Buffer) Len() int {
	return len(b.complete) + len(b.incomplete)
}

// Consume clears the buffer of the current line, and feeds any remaining
// buffered data back into WriteString.
func (b *Buffer) Reset() error {
	incomplete := b.incomplete
	*b = Buffer{}
	_, err := b.WriteString(incomplete)
	return err
}

// Ready returns true iff the buffer contains a complete line.
func (b *Buffer) Ready() bool {
	return b.complete != ""
}

// Tokens returns the tokens added to the buffer.
func (b *Buffer) Tokens() []token.Token {
	return b.tokens
}

// TODO revise comment
// String returns the contents of the buffer as a string. If the Buffer is a
// nil pointer, it returns "<nil>".
func (b *Buffer) String() string {
	return b.complete
}
