package zone

type TokenType int

const (
	BracketCloseToken TokenType = iota
	BracketOpenToken
	CommentToken
	RecordToken
	SpaceToken
	CharToken
)

func (t TokenType) String() string {
	return []string{")", "(", ";", "", " ", ""}[t]
}

type Token struct {
	Type   TokenType
	Data   string
	Tokens []Token
}

func NewToken(t TokenType) Token {
	return Token{
		Type: t,
	}
}

func NewCharToken(c byte) Token {
	return Token{
		Type: CharToken,
		Data: string(c),
	}
}

func (t *Token) AddTokens(tokens []Token) {
	t.Tokens = append(t.Tokens, tokens...)
}

func (t *Token) String() string {
	switch t.Type {
	case BracketCloseToken, BracketOpenToken, CommentToken, SpaceToken:
		return t.Type.String()
	case RecordToken:
		s := ""
		for _, token := range t.Tokens {
			s += token.String()
		}
		return s
	default:
		return t.Data
	}
}
