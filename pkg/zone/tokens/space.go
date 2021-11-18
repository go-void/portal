package tokens

type Space struct {
	tokens []Token
}

func (t *Space) String() string {
	return t.tokens[0].String()
}

func (t *Space) Tokens() []Token {
	return t.tokens
}

func (t *Space) Add(_ []Token) {
	t.tokens = append(t.tokens, &Char{
		Data: " ",
	})
}

func (t *Space) Type() int {
	return TypeSpace
}
