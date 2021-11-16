package zone

const (
	TypeChar int = iota
	TypeSpace
	TypeComment
	TypeRecord
)

type Token interface {
	Tokens() []Token
	Data() string
	Add([]Token)
	Type() int
}

func NewToken(t int) Token {
	switch t {
	case TypeChar:
		return new(Char)
	case TypeSpace:
		return new(Space)
	case TypeComment:
		return new(Comment)
	case TypeRecord:
		return new(Record)
	}
	return nil
}

type Char struct {
	data string
}

func (t *Char) Data() string {
	return t.data
}

func (t *Char) Tokens() []Token {
	return nil
}

func (t *Char) Add(_ []Token) {
	return
}

func (t *Char) Type() int {
	return TypeChar
}

type Space struct {
	tokens []Token
}

func (t *Space) Data() string {
	return t.tokens[0].Data()
}

func (t *Space) Tokens() []Token {
	return t.tokens
}

func (t *Space) Add(_ []Token) {
	t.tokens = append(t.tokens, &Char{
		data: " ",
	})
}

func (t *Space) Type() int {
	return TypeSpace
}

type Comment struct {
	tokens []Token
}

func (t *Comment) Data() string {
	var s string
	for _, token := range t.tokens {
		s += token.Data()
	}
	return s
}

func (t *Comment) Tokens() []Token {
	return t.tokens
}

func (t *Comment) Add(tokens []Token) {
	t.tokens = append(t.tokens, tokens...)
}

func (t *Comment) Type() int {
	return TypeComment
}

type Record struct {
	tokens []Token
}

func (t *Record) Data() string {
	var s string
	for _, token := range t.tokens {
		s += token.Data()
	}
	return s
}

func (t *Record) Tokens() []Token {
	return t.tokens
}

func (t *Record) Add(tokens []Token) {
	t.tokens = append(t.tokens, tokens...)
}

func (t *Record) Type() int {
	return TypeRecord
}
