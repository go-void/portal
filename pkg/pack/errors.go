package pack

type OverflowError string

func (e OverflowError) Error() string {
	return string(e)
}

const (
	ErrOverflowUnpackUint64 = OverflowError("offset overflow unpacking uint64")
	ErrOverflowUnpackUint32 = OverflowError("offset overflow unpacking uint32")
	ErrOverflowUnpackUint16 = OverflowError("offset overflow unpacking uint16")
	ErrOverflowUnpackUint8  = OverflowError("offset overflow unpacking uint8")

	ErrOverflowUnpackIPv4 = OverflowError("offset overflow unpacking IPv4 address")
	ErrOverflowUnpackIPv6 = OverflowError("offset overflow unpacking IPv6 address")

	ErrOverflowUnpackString = OverflowError("offset overflow unpacking character string")
	ErrOverflowUnpackName   = OverflowError("offset overflow unpacking domain name")
)
