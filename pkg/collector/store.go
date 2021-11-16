package collector

type CollectorStore interface {
	CreateEntries([]Entry) error
}

type DefaultCollectorStore struct {
}

func NewDefaultStore() *DefaultCollectorStore {
	return &DefaultCollectorStore{}
}

func (s *DefaultCollectorStore) CreateEntries(_ []Entry) error {
	return nil
}
