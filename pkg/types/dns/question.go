package dns

import "go.uber.org/zap/zapcore"

// Questions holds a DNS question. The RFC allows multiple questions per
// message, but most DNS servers only accpet one and multiple questions
// often result in errors.
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

// MarshalLogObject marshals a DNS question as a zap log object
func (q Question) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", q.Name)
	enc.AddUint16("type", q.Type)
	enc.AddUint16("class", q.Class)
	return nil
}
