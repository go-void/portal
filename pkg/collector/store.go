package collector

// CollectorStore defines an interface to flush in-memory entries to persist via a backend
type CollectorStore interface {
	// StoreEntries stores a slice of entries
	StoreEntries([]Entry) error

	// Prepare prepares the backend (if needed)
	Prepare() error
}

// DefaultCollectorStore discards flushed entries
type DefaultCollectorStore struct {
}

func NewDefaultStore() *DefaultCollectorStore {
	return &DefaultCollectorStore{}
}

func (s *DefaultCollectorStore) StoreEntries(_ []Entry) error {
	return nil
}

func (s *DefaultCollectorStore) Prepare() error {
	return nil
}
