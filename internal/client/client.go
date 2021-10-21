package client

type Client interface {
	Query(string, uint16, uint16, QueryOptions) (interface{}, error)
}

type DefaultClient struct {
}

type QueryOptions struct {
}

func NewDefaultClient() *DefaultClient {
	return &DefaultClient{}
}

func (c *DefaultClient) Query(name string, class, t uint16, options QueryOptions) (interface{}, error) {
	// Construct query / question message with options

	// Send message

	// Wait for answer

	// Return answer or error
	return nil, nil
}
