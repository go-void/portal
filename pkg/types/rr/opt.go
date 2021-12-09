package rr

import (
	"fmt"

	"github.com/go-void/portal/pkg/pack"
	"github.com/go-void/portal/pkg/types/edns"
)

// See https://datatracker.ietf.org/doc/html/rfc6891#section-6.1
// Many OPT header fields have a special usa case:
// - Class: UDP payload size. RFC 6891 romoved the 512 byte size limit.
//   This field can set the maximum size in bytes (e.g. 4096)
// - TTL: This 32 bit field is split up in:
//   - 1 octet extended RCODEs
//   - 1 octet EDNS version
//   - DO bit
//   - 15 reserved bits

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
	return fmt.Sprintf("%v %v", rr.H, rr.Options)
}

func (rr *OPT) Len() uint16 {
	var len = uint16(0)
	for _, o := range rr.Options {
		len += o.Len()
	}
	return len
}

func (rr *OPT) Unpack(data []byte, offset int) (int, error) {
	options, offset, err := pack.UnpackEDNSOptions(data, offset, rr.H.RDLength)
	if err != nil {
		return offset, err
	}
	rr.Options = options
	return offset, nil
}

func (rr *OPT) Pack(buf []byte, offset int) (int, error) {
	return offset + 4, nil
}
