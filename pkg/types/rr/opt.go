package rr

import (
	"github.com/go-void/portal/pkg/types/edns"
)

// See https://datatracker.ietf.org/doc/html/rfc6891#section-6.1
// The 32 bit TTL field in the OPT RR header has a special use case:
// +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
// |         EXTENDED-RCODE        |            VERSION            |
// +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
// | DO|                           Z                               |
// +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+

type OPT struct {
	H       Header
	Options []edns.Option
}

func (rr *OPT) Header() *Header {
	return &rr.H
}

func (rr *OPT) SetHeader(header Header) {
	rr.H = header
}

func (rr *OPT) SetData(data ...interface{}) error {
	return nil
}

func (rr *OPT) String() string {
	return ""
}

func (rr *OPT) Len() uint16 {
	return 0
}

func (rr *OPT) Unpack(data []byte, offset int) (int, error) {
	return offset, nil
}

func (rr *OPT) Pack(buf []byte, offset int) (int, error) {
	return offset + 4, nil
}
