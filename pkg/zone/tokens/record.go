package tokens

type Record struct {
	tokens []Token
}

func (t *Record) String() string {
	var s string
	for _, token := range t.tokens {
		s += token.String()
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
