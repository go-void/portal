package pack

import (
	"github.com/go-void/portal/internal/labels"
	"github.com/go-void/portal/internal/types/dns"
	"github.com/go-void/portal/internal/types/rr"
)

// Packer packs DNS messages from a struct into the
// wire format
type Packer interface {
	// Packs packs a single DNS message by converting
	// the provided message to the wire format
	Pack(dns.Message) ([]byte, error)

	// PackHeader packs header data by converting
	// the provided header to the wire format
	PackHeader(dns.MessageHeader, []byte, int) (int, error)

	// PackQuestion packs a question by converting
	// the provided question to the wire format
	PackQuestion(dns.Question, []byte, int) (int, error)

	PackName(string, []byte, int) (int, error)

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
	PackRRHeader(*rr.RRHeader, []byte, int) (int, error)
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

	_, err = p.PackRRList(message.Additional, buf, offset)
	return buf[:offset], err
}

// PackHeader packs header data by converting the provided header to the wire format
func (p *DefaultPacker) PackHeader(header dns.MessageHeader, buf []byte, offset int) (int, error) {
	rh := header.ToRaw()

	offset, err := PackUint16(rh.ID, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(rh.Flags, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(rh.QDCount, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(rh.ANCount, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(rh.NSCount, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(rh.ARCount, buf, offset)
	return offset, err
}

// PackQuestion packs a question by converting the provided question to the wire format
func (p *DefaultPacker) PackQuestion(question dns.Question, buf []byte, offset int) (int, error) {
	offset, err := p.PackName(question.Name, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(question.Type, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(question.Class, buf, offset)
	return offset, err
}

func (p *DefaultPacker) PackName(name string, buf []byte, offset int) (int, error) {
	labels := labels.FromBottom(name)

	for i := 0; i < len(labels); i++ {
		label := labels[i]
		switch label {
		case "", ".":
			buf[offset] = 0x0
			offset++
		default:
			buf[offset] = uint8(len(label))
			offset++

			for l := 0; l < len(label); l++ {
				buf[offset] = label[l]
				offset++
			}
		}
	}

	return offset, nil
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
func (p *DefaultPacker) PackRRHeader(header *rr.RRHeader, buf []byte, offset int) (int, error) {
	offset, err := p.PackName(header.Name, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(header.Type, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(header.Class, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint32(header.TTL, buf, offset)
	if err != nil {
		return offset, err
	}

	offset, err = PackUint16(header.RDLength, buf, offset)
	return offset, err
}
