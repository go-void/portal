package cache

import "github.com/go-void/portal/internal/types/rr"

type Cache interface {
	Get(string, uint16, uint16) (rr.RR, bool)

	Set(string, string)

	Len() int
}
