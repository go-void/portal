package resolver

import (
	"github.com/go-void/portal/pkg/types/dns"
	"github.com/go-void/portal/pkg/types/rr"
)

type Result struct {
	Answer     []rr.RR
	Authority  []rr.RR
	Additional []rr.RR
}

func NewResult(message *dns.Message) Result {
	return Result{message.Answer, message.Authority, message.Additional}
}
