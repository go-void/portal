package dns

// Questions holds a DNS question. The RFC allows multiple questions per
// message, but most DNS servers only accpet one and multiple questions
// often result in errors.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
type Question struct {
	Name  string
	Type  uint16
	Class uint16
}
