package zone

import (
	"bufio"
	"errors"
	"io"
)

var (
	ErrExtraClosingBracket = errors.New("extra closing bracket")
)

type Tokenizer struct {
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

	// Buff holds a series of tokens until they are consumed
	Buff []Token

	// Error is non-nil if the lexer encountered an error
	// along the way of tokenizing the input
	Error error

	// Tokens hold a series of tokens of the final result
	Tokens []Token
}

func NewTokenizer(r io.Reader) *Tokenizer {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = bufio.NewReaderSize(r, 1024)
	}

	return &Tokenizer{
		Reader: br,
	}
}

func (t *Tokenizer) Parse() ([]Token, error) {
	for {
		ok := t.next()
		if !ok {
			break
		}
	}

	return t.Tokens, t.Error
}

func (t *Tokenizer) read() bool {
	c, err := t.Reader.ReadByte()
	// TODO (Techassi): Handle EOF "error"
	if err != nil {
		t.Error = err
		return false
	}

	t.Current = c
	return true
}

func (t *Tokenizer) next() bool {
	if t.Error != nil {
		return false
	}

	for t.read() {
		switch t.Current {
		case ';':
			if t.Quoted || t.Escaped {
				t.Escaped = false
				t.Buff = append(t.Buff, NewToken(CommentToken))
				continue
			}
			t.InComment = true
		case '\n', '\r':
			if t.Brackets > 0 {
				continue
			}

			if len(t.Buff) == 0 {
				continue
			}

			if t.InComment {
				token := NewToken(CommentToken)
				token.AddTokens(t.Buff)
				t.Tokens = append(t.Tokens, token)

				t.Buff = []Token{}
				t.InComment = false
				continue
			}

			token := NewToken(RecordToken)
			token.AddTokens(t.Buff)
			t.Tokens = append(t.Tokens, token)

			t.Buff = []Token{}
		case '"':
			if t.Quoted && !t.Escaped {
				t.Quoted = false
			}
		case '(', ')':
			if t.Escaped || t.Quoted {
				t.Escaped = false
			}

			if t.Current == '(' {
				t.Brackets++
				token := NewToken(BracketOpenToken)
				t.Buff = append(t.Buff, token)
				continue
			}

			if t.Current == ')' {
				t.Brackets--
			}

			if t.Brackets < 0 {
				t.Error = ErrExtraClosingBracket
				return false
			}

			token := NewToken(BracketCloseToken)
			t.Buff = append(t.Buff, token)
		case '\\':
			t.Escaped = true
		case ' ', '\t':
			if t.Escaped || t.Quoted || t.InComment {
				t.Escaped = false
			}

			token := NewToken(SpaceToken)
			t.Buff = append(t.Buff, token)
		default:
			t.Buff = append(t.Buff, NewCharToken(t.Current))
		}
	}

	return false
}
