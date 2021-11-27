package filter

type Store interface {
	// Add adds a new filter to the store
	Add(string, Filter) error

	// Get returns a filter by name
	Get(string) (*Filter, error)
}

type DefaultStore struct {
}

func NewDefault() Store {
	return &DefaultStore{}
}

func (s *DefaultStore) Add(name string, f Filter) error {
	// Code ...
	return nil
}

func (s *DefaultStore) Get(name string) (*Filter, error) {
	// Code ...
	return nil, nil
}
