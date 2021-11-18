package tokens

const (
	TypeChar int = iota
	TypeSpace
	TypeComment
	TypeRecord
	TypeBracketOpen
	TypeBracketClose
)

type Token interface {
	Tokens() []Token
	String() string
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
	case TypeBracketOpen:
		return new(Record)
	}
	return nil
}
