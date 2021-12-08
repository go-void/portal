package edns

import "errors"

var (
	ErrNoSuchOption = errors.New("no such option")
)

type Option interface {
	// Code returns the option code
	Code() uint16

	// Len returns the option length
	Len() uint16

	// Unpack unpacks the option data
	Unpack([]byte, int, uint16) (int, error)

	// Pack packs the option data
	Pack([]byte, int) (int, error)
}

func New(code uint16) (Option, error) {
	create, ok := codeMap[code]
	if !ok {
		return nil, ErrNoSuchOption
	}

	return create(), nil
}
