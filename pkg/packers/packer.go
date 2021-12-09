package packers

import (
	"github.com/go-void/portal/pkg/pack"
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

// Packer packs DNS messages from a struct into the
// wire format
type Packer interface {
	// Packs packs a single DNS message by converting
	// the provided message to the wire format
	Pack(dns.Message) ([]byte, error)

	// PackHeader packs header data by converting
	// the provided header to the wire format
	PackHeader(dns.Header, []byte, int) (int, error)

	// PackQuestion packs a question by converting
	// the provided question to the wire format
	PackQuestion(dns.Question, []byte, int) (int, error)

	// PackRRList packs a slice of resource records
	// by converting the provided records to the
	// wire format
	PackRRList([]rr.RR, []byte, int) (int, error)

	// PackRR packs a single resource record by
	// converting the provided data to the wire
	// format
	PackRR(rr.RR, []byte, int) (int, error)

	// PackRRHeader packs a resource record header
	// by converting the provided data to the
	// wire format
	PackRRHeader(*rr.Header, []byte, int) (int, error)
}

// DefaultPacker is the default packer implementation
// which follows the specs RFC 1034 and 1035
type DefaultPacker struct {
}

func NewDefaultPacker() Packer {
	return &DefaultPacker{}
}

// Packs packs a single DNS message by converting the provided
// message to the wire format
func (p *DefaultPacker) Pack(message dns.Message) ([]byte, error) {
	// FIXME (Techassi): Figure out how we can pre-allocate the buf with the correct length / size
	var buf = make([]byte, 256*4)

	offset, err := p.PackHeader(message.Header, buf, 0)
	if err != nil {
		return buf, err
	}

	offset, err = p.PackQuestion(message.Question[0], buf, offset)
	if err != nil {
		return buf, err
	}

	offset, err = p.PackRRList(message.Answer, buf, offset)
	if err != nil {
		return buf, err
	}

	offset, err = p.PackRRList(message.Authority, buf, offset)
	if err != nil {
		return buf, err
	}

	offset, err = p.PackRRList(message.Additional, buf, offset)
	return buf[:offset], err
}

// PackHeader packs header data by converting the provided header to the wire format
func (p *DefaultPacker) PackHeader(header dns.Header, buf []byte, offset int) (int, error) {
	rh := header.ToRaw()

	offset, err := pack.PackUint16(rh.ID, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(rh.Flags, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(rh.QDCount, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(rh.ANCount, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(rh.NSCount, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(rh.ARCount, buf, offset)
	return offset, err
}

// PackQuestion packs a question by converting the provided question to the wire format
func (p *DefaultPacker) PackQuestion(question dns.Question, buf []byte, offset int) (int, error) {
	offset, err := pack.PackDomainName(question.Name, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(question.Type, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(question.Class, buf, offset)
	return offset, err
}

// PackRRList packs a slice of resource records by converting the provided records to
// the wire format
func (p *DefaultPacker) PackRRList(rrs []rr.RR, buf []byte, offset int) (int, error) {
	if len(rrs) == 0 {
		return offset, nil
	}

	var err error
	for _, rr := range rrs {
		offset, err = p.PackRR(rr, buf, offset)
		if err != nil {
			return offset, err
		}
	}
	return offset, nil
}

// PackRR packs a single resource record by converting the provided
// data to the wire format
func (p *DefaultPacker) PackRR(rr rr.RR, buf []byte, offset int) (int, error) {
	offset, err := p.PackRRHeader(rr.Header(), buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = rr.Pack(buf, offset)
	return offset, err
}

// PackRRHeader packs a resource record header by converting the provided data
// to the wire format
func (p *DefaultPacker) PackRRHeader(header *rr.Header, buf []byte, offset int) (int, error) {
	offset, err := pack.PackDomainName(header.Name, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(header.Type, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(header.Class, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint32(header.TTL, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = pack.PackUint16(header.RDLength, buf, offset)
	return offset, err
}
