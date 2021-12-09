package edns

import (
	"encoding/hex"
	"fmt"
)

type Cookie struct {
	Client string
	Server string
}

// Code returns the option code
func (o *Cookie) Code() uint16 {
	return CodeCOOKIE
}

// Len returns the option length
func (o *Cookie) Len() uint16 {
	return 8 + uint16(len(o.Server))
}

func (o *Cookie) String() string {
	return fmt.Sprintf("Client: %s, Server: %s", o.Client, o.Server)
}

// Unpack unpacks the option data
func (o *Cookie) Unpack(data []byte, offset int, length uint16) (int, error) {
	client := data[offset : offset+8]
	o.Client = hex.EncodeToString(client)
	offset += 8

	// If the length is only 8 (octects) the client only send
	// a client cookie and the server cookie is empty. Return
	// early
	if length == 8 {
		return offset, nil
	}

	server := data[offset : offset+int(length-8)]
	o.Server = hex.EncodeToString(server)
	offset += int(length) - 8

	return offset, nil
}

// Pack packs the option data
func (o *Cookie) Pack(_ []byte, _ int) (int, error) {
	panic("not implemented") // TODO: Implement
}
