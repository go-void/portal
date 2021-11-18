package zone

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/go-void/portal/pkg/types/rr"
	"github.com/go-void/portal/pkg/zone/tokens"
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
	Buff    []tokens.Token

	// Error is non-nil if the lexer encountered
	// an error along the way of tokenizing the
	// input
	Error error

	Tokens []tokens.Token
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

func (l *Parser) Parse() ([]tokens.Token, error) {
	for {
		ok := l.next()
		if !ok {
			break
		}
	}
	fmt.Println(len(l.Tokens))
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

	// TODO (Techassi): Handle zone files which don't end with a empty newline. Currently we discard the buff. This
	// makes it neccesary to insert a newline at the end of file.

	for l.read() {
		switch l.Current {
		case ';':
			if l.Quoted || l.Escaped {
				l.Escaped = false
				l.Buff = append(l.Buff, &tokens.Char{
					Data: ";",
				})
				continue
			}
			l.InComment = true
		case '\n', '\r':
			if l.Brackets > 0 {
				continue
			}

			if len(l.Buff) == 0 {
				continue
			}

			if l.InComment {
				t := tokens.NewToken(tokens.TypeComment)
				t.Add(l.Buff)
				l.Tokens = append(l.Tokens, t)

				l.Buff = []tokens.Token{}
				l.InComment = false
				continue
			}

			t := tokens.NewToken(tokens.TypeRecord)
			t.Add(l.Buff)
			l.Tokens = append(l.Tokens, t)

			l.Buff = []tokens.Token{}
		case '"':
			if l.Quoted && !l.Escaped {
				l.Quoted = false
			}
		case '(', ')':
			if l.Escaped || l.Quoted {
				l.Escaped = false
			}

			if l.Current == '(' {
				l.Brackets++
				t := tokens.NewToken(tokens.TypeBracketOpen)
				l.Buff = append(l.Buff, t)
				continue
			}

			if l.Current == ')' {
				l.Brackets--
			}

			if l.Brackets < 0 {
				l.Error = ErrExtraClosingBracket
				return false
			}

			t := tokens.NewToken(tokens.TypeBracketClose)
			l.Buff = append(l.Buff, t)
		case '\\':
			l.Escaped = true
		case ' ', '\t':
			if l.Escaped || l.Quoted || l.InComment {
				l.Escaped = false
			}

			t := tokens.NewToken(tokens.TypeSpace)
			t.Add(nil)
			l.Buff = append(l.Buff, t)
		default:
			l.Buff = append(l.Buff, &tokens.Char{
				Data: string(l.Current),
			})
		}
	}

	return false
}
