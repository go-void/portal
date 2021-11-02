package zone

import (
	"bufio"
	"errors"
	"io"

	"github.com/go-void/portal/internal/types/rr"
)

var (
	ErrExtraClosingBracket = errors.New("extra closing bracket")
)

type ZoneParser interface {
	Read() bool
	Next() bool
	Peek() ZoneParser
}

type Parser struct {
	// Input byte reader
	Reader io.ByteReader

	// Context
	InComment bool
	Escaped   bool
	Quoted    bool
	EOL       bool
	Brackets  int

	// Current stores the current looked at byte
	Current byte

	// Positional states
	Column int
	Line   int

	// Error is non-nil if the lexer encountered
	// an error along the way of tokenizing the
	// input
	Error error

	RRs []rr.RR
}

func NewParser(r io.Reader) *Parser {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = bufio.NewReaderSize(r, 1024)
	}

	return &Parser{
		Reader: br,
	}
}

func (l *Parser) Read() bool {
	c, err := l.Reader.ReadByte()
	if err != nil {
		l.Error = err
		return false
	}

	if c == '\n' {
		l.EOL = true
		l.Line++
	} else {
		l.Column++
	}

	l.Current = c
	return true
}

func (l *Parser) Next() bool {
	if l.Error != nil {
		return false
	}

	for l.Read() {
		switch l.Current {
		case ';':
			if l.Quoted || l.Escaped {
				l.Escaped = false
				continue
			}
			l.InComment = true
		case '\n':
			l.InComment = false

			if l.Brackets > 0 {
				continue
			}

			// We have a complete RR
		case '"':
			if l.Quoted && !l.Escaped {
				l.Quoted = false
			}
		case '(', ')':
			if l.Escaped || l.Quoted {
				l.Escaped = false
				continue
			}

			if l.Current == '(' {
				l.Brackets++
				continue
			}

			if l.Current == ')' {
				l.Brackets--

				if l.Brackets < 0 {
					l.Error = ErrExtraClosingBracket
					return false
				}
			}
		case '\\':
			l.Escaped = true
		case ' ', '\t':
			if l.Escaped || l.Quoted || l.InComment {
				l.Escaped = false
				continue
			}
			// Regular RR data
		}
	}

	return true
}

func (l *Parser) Peek() Parser {
	panic("not implemented") // TODO: Implement
}
