package tokens

type Char struct {
	Data string
}

func (t *Char) String() string {
	return t.Data
}

func (t *Char) Tokens() []Token {
	return nil
}

func (t *Char) Add(_ []Token) {}

func (t *Char) Type() int {
	return TypeChar
}
