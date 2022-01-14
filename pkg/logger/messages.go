package logger

// General log messages
const (
	ErrUnpackDNSHeader  = "failed to unpack DNS message header"
	ErrUnpackDNSMessage = "failed to unpack DNS message"
	ErrPackDNSMessage   = "failed to pack DNS message"
	ErrAcceptMessage    = "did not accept DNS message"
	ErrHandleRequest    = "failed to handle incoming DNS request"
	ErrResolverLookup   = "failed to lookup domain name via resolver"
)

// UDP related log messages
const (
	ErrUDPRead  = "failed to read UDP datagram"
	ErrUDPWrite = "failed to write UDP datagram"
)

// TCP related log messages
const (
	ErrTCPAccept     = "failed to accept TCP conn"
	ErrTCPClose      = "failed to close TCP conn"
	ErrTCPRead       = "failed to read TCP packet"
	ErrTCPWrite      = "failed to write TCP packet"
	ErrTCPWriteClose = "failed to write packet or close TCP conn"
)

// Filter related log messages
const (
	DebugNoSuchFilter = "no matching filter found"
)
