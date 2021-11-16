package zone

import (
	"bufio"
	"errors"
	"io"

	"github.com/go-void/portal/pkg/types/rr"
)

var (
	ErrExtraClosingBracket = errors.New("extra closing bracket")
)

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
	Buff    []Token

	// Error is non-nil if the lexer encountered
	// an error along the way of tokenizing the
	// input
	Error error

	Tokens []Token
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

func (l *Parser) Parse() ([]Token, error) {
	for {
		ok := l.next()
		if !ok {
			break
		}
	}
	return l.Tokens, l.Error
}

func (l *Parser) Records() ([]rr.RR, error) {

	return nil, nil
}

func (l *Parser) read() bool {
	c, err := l.Reader.ReadByte()
	// TODO (Techassi): Handle EOF "error"
	if err != nil {
		l.Error = err
		return false
	}

	l.Current = c
	return true
}

func (l *Parser) next() bool {
	if l.Error != nil {
		return false
	}

	for l.read() {
		switch l.Current {
		case ';':
			if l.Quoted || l.Escaped {
				l.Escaped = false
				l.Buff = append(l.Buff, &Char{
					data: ";",
				})
				continue
			}
			l.InComment = true
		case '\n':
			if l.Brackets > 0 {
				continue
			}

			if l.InComment {
				t := NewToken(TypeComment)
				t.Add(l.Buff)
				l.Tokens = append(l.Tokens, t)

				l.Buff = []Token{}
				l.InComment = false
				continue
			}

			t := NewToken(TypeRecord)
			t.Add(l.Buff)
			l.Tokens = append(l.Tokens, t)

			l.Buff = []Token{}
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
			}

			t := NewToken(TypeSpace)
			t.Add(nil)
			l.Buff = append(l.Buff, t)
		default:
			l.Buff = append(l.Buff, &Char{
				data: string(l.Current),
			})
		}
	}

	return true
}
