package pack

import (
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
	PackHeader(dns.MessageHeader) ([]byte, int, error)

	// PackQuestion packs a question by converting
	// the provided question to the wire format
	PackpQuestion(dns.Question) ([]byte, int, error)

	// PackRRList packs a slice of resource records
	// by converting the provided records to the
	// wire format
	PackRRList([]rr.RR) ([]byte, int, error)

	// PackRR packs a single resource record by
	// converting the provided data to the wire
	// format
	PackRR(rr.RR) ([]byte, int, error)

	// PackRRHeader packs a resource record header
	// by converting the provided data to the
	// wire format
	PackRRHeader(rr.RRHeader) ([]byte, int, error)
}

// DefaultPacker is the default packer implementation
// which follows the specs RFC 1034 and 1035
type DefaultPacker struct {
}

func NewDefaultPackper() Packer {
	return &DefaultPacker{}
}

// Packs packs a single DNS message by converting the provided
// message to the wire format
func (p *DefaultPacker) Pack(_ dns.Message) ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

// PackHeader packs header data by converting the provided header to the wire
// format
func (p *DefaultPacker) PackHeader(_ dns.MessageHeader) ([]byte, int, error) {
	panic("not implemented") // TODO: Implement
}

// PackQuestion packs a question by converting the provided question to the
// wire format
func (p *DefaultPacker) PackpQuestion(_ dns.Question) ([]byte, int, error) {
	panic("not implemented") // TODO: Implement
}

// PackRRList packs a slice of resource records by converting the provided
// records to the wire format
func (p *DefaultPacker) PackRRList(_ []rr.RR) ([]byte, int, error) {
	panic("not implemented") // TODO: Implement
}

// PackRR packs a single resource record by converting the provided
// data to the wire format
func (p *DefaultPacker) PackRR(_ rr.RR) ([]byte, int, error) {
	panic("not implemented") // TODO: Implement
}

// PackRRHeader packs a resource record header by converting the provided data
// to the wire format
func (p *DefaultPacker) PackRRHeader(_ rr.RRHeader) ([]byte, int, error) {
	panic("not implemented") // TODO: Implement
}
