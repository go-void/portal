package zone

import (
	"bufio"
	"io"
	"strings"
)

type ZoneLexer interface {
	Read() bool
	Next() bool
	Peek() ZoneLexer
}

type Lexer struct {
	// Input byte reader
	Reader io.ByteReader

	// Ctx stores the current context of the lexer
	// which keeps track if the lexer is inside a
	// quoted string for example
	Ctx LexerContext

	// Current stores the current looked at byte
	Current byte

	// Positional states
	Column int
	Line   int

	// Error is non-nil if the lexer encountered
	// an error along the way of tokenizing the
	// input
	Error error

	// Multiple result value string builders
	Result  strings.Builder
	Comment strings.Builder

	Tokens []Token

	BlankCount int
}

type LexerContext struct {
	InComment bool
	Escaped   bool
	Quoted    bool
	EOL       bool
}

func NewLexer(r io.Reader) *Lexer {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = bufio.NewReaderSize(r, 1024)
	}

	return &Lexer{
		Reader: br,
		Ctx:    LexerContext{},
	}
}

func (l *Lexer) Read() bool {
	c, err := l.Reader.ReadByte()
	if err != nil {
		l.Error = err
		return false
	}

	if c == '\n' {
		l.Ctx.EOL = true
		l.Line++
	} else {
		l.Column++
	}

	l.Current = c
	return true
}

func (l *Lexer) Next() bool {
	if l.Error != nil {
		return false
	}

	for l.Read() {
		switch l.Current {
		case ' ', '\t':
			if l.Ctx.InComment {
				break
			}

			if l.BlankCount == 0 {
				l.Tokens = append(l.Tokens, BlankT)
			}
			l.BlankCount++

			// Usually whitespaces and tabs delimit fields of one entry.
			// If we are inside a quoted or escaped string this is valid
			// if l.Ctx.Escaped || l.Ctx.Quoted {
			// 	l.Result.WriteByte(l.Current)
			// 	l.unEscape()
			// 	break
			// }

			// Whitespaces and tabs are also allowed in comments
			// if l.Ctx.InComment {
			// 	l.Comment.WriteByte(l.Current)
			// 	break
			// }

			// TODO (Techassi): Handle other cases of whitespaces and tabs
		case ';':
			l.inComment()

			// Usually semicolons indicate a comment. If we are inside a
			// quoted or escaped string this is valid
			// if l.Ctx.Escaped || l.Ctx.Quoted {
			// 	l.Result.WriteByte(l.Current)
			// 	l.unEscape()
			// 	break
			// }

			// l.inComment()

			// l.Comment.WriteByte(';')
			// TODO (Techassi): Handle comment count and newlines in comments
		case '\n':
			// Handle newlines
			if l.Ctx.InComment {
				l.outComment()
			}

			l.Tokens = append(l.Tokens, NewlineT)
		case '\r':
			break
			// Carriage returns are only allowed in quoted string or escaped.
			// Everything else is ignored / skipped
			// l.unEscape()

			// if l.Ctx.Quoted {
			// 	l.Result.WriteByte(l.Current)
			// }
		case '"':
			// TODO (Techassi): Handle quoted strings
			if l.Ctx.InComment {
				break
			}
			l.Tokens = append(l.Tokens, QuoteT)
		case '\\':
			// TODO (Techassi): Handle escaped strings
		case '(', ')':
			// TODO (Techassi): Handle brackets
			if l.Ctx.InComment {
				break
			}

			if l.Current == '(' {
				l.Tokens = append(l.Tokens, OpenBracketT)
			} else {
				l.Tokens = append(l.Tokens, CloseBracketT)
			}
		default:
			// TODO (Techassi): Handle default
		}
	}

	return true
}

func (l *Lexer) Peek() Lexer {
	panic("not implemented") // TODO: Implement
}

// func (l *ZoneLexer) escape() {
// 	l.Ctx.Escaped = true
// }

// func (l *ZoneLexer) unEscape() {
// 	l.Ctx.Escaped = false
// }

func (l *Lexer) inComment() {
	l.Ctx.InComment = true
}

func (l *Lexer) outComment() {
	l.Ctx.InComment = false
}
