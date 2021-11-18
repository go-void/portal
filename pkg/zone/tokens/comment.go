package tokens

type Comment struct {
	tokens []Token
}

func (t *Comment) String() string {
	var s string
	for _, token := range t.tokens {
		s += token.String()
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
