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
	ErrOverflowPackUint32   = OverflowError("offset overflow packing uint32")
	ErrOverflowPackUint16   = OverflowError("offset overflow packing uint16")
	ErrOverflowPackUint8    = OverflowError("offset overflow packing uint8")

	ErrOverflowUnpackIPv4 = OverflowError("offset overflow unpacking IPv4 address")
	ErrOverflowUnpackIPv6 = OverflowError("offset overflow unpacking IPv6 address")
	ErrOverflowPackIPv4   = OverflowError("offset overflow packing IPv4 address")
	ErrOverfloPpackIPv6   = OverflowError("offset overflow packing IPv6 address")

	ErrOverflowUnpackString = OverflowError("offset overflow unpacking character string")
	ErrOverflowUnpackName   = OverflowError("offset overflow unpacking domain name")
)
