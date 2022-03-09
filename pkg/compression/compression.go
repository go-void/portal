package compression

import "errors"

var (
	ErrNoSuchCompressionPointer = errors.New("compression: no such pointer")
	ErrNoSuchCompressionName    = errors.New("compression: no such name")
)

// Map holds data for compressing DNS messages.
type Map struct {
	m map[string]int
}

func New() Map {
	return Map{
		m: make(map[string]int),
	}
}

func (c Map) Set(name string, ptr int) {
	_, ok := c.m[name]
	if ok {
		return
	}

	c.m[name] = ptr
}

func (c Map) Get(name string) int {
	return 0
}
