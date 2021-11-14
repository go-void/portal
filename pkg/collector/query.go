package collector

import (
	"net"
	"time"

	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type Entry struct {
	QueryTime time.Duration
	Question  dns.Question
	Answer    rr.RR
	Result    string
	ClientIP  net.IP
	Filtered  bool
	Cached    bool
}
