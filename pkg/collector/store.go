package collector

// Store defines an interface to flush in-memory entries to persist via a backend
type Store interface {
	// StoreEntries stores a slice of entries
	StoreEntries([]Entry) error

	// Prepare prepares the backend (if needed)
	Prepare() error
}

// DefaultStore discards flushed entries
type DefaultStore struct {
}

func NewStore() *DefaultStore {
	return &DefaultStore{}
}

func (s *DefaultStore) StoreEntries(_ []Entry) error {
	return nil
}

func (s *DefaultStore) Prepare() error {
	return nil
}
