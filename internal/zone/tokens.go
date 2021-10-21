package zone

type Token int

const (
	EOFT Token = iota
	BlankT
	StringT
	NewlineT
	QuoteT
	OpenBracketT
	CloseBracketT
)
