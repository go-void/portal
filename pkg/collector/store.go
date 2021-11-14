package collector

type CollectorStore interface {
	CreateFilterEntry() error

	CreateQueryEntry() error
}

type DefaultCollectorStore struct {
}

func NewDefaultStore() *DefaultCollectorStore {
	return &DefaultCollectorStore{}
}

func (s *DefaultCollectorStore) CreateFilterEntry() error {
	// Code ...
	return nil
}

func (s *DefaultCollectorStore) CreateQueryEntry() error {
	// Code ...
	return nil
}
