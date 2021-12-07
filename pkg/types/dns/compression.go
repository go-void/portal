package dns

import "errors"

var (
	ErrNoSuchCompressionPointer = errors.New("compression: no such pointer")
	ErrNoSuchCompressionName    = errors.New("compression: no such name")
)

// CompressionMap holds data for (de)compressing DNS messages.
type CompressionMap struct {
	ptrToName map[int]string
	nameToPtr map[string]int
}

// NewCompressionMap returns a new empty compression map
func NewCompressionMap() CompressionMap {
	return CompressionMap{
		ptrToName: make(map[int]string),
		nameToPtr: make(map[string]int),
	}
}

func (m *CompressionMap) ReturnOrAddPointer(name string, ptr int) int {
	if p, ok := m.nameToPtr[name]; ok {
		return p
	}
	m.nameToPtr[name] = ptr
	return ptr
}

func (m *CompressionMap) ReturnOrAddName(ptr int, name string) string {
	if n, ok := m.ptrToName[ptr]; ok {
		return n
	}
	m.ptrToName[ptr] = name
	return name
}

// GetPointer returns a pointer which references the given domain name
func (m *CompressionMap) GetPointer(name string) (int, error) {
	if ptr, ok := m.nameToPtr[name]; ok {
		return ptr, nil
	}
	return -1, ErrNoSuchCompressionPointer
}

// GetName returns a name which is referenced by the given pointer
func (m *CompressionMap) GetName(ptr int) (string, error) {
	if name, ok := m.ptrToName[ptr]; ok {
		return name, nil
	}
	return "", ErrNoSuchCompressionName
}
