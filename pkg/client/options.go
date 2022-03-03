package client

import "time"

type OptionFunc func(*Client) error

func WithNetwork(network string) OptionFunc {
	return func(c *Client) error {
		c.network = network
		return nil
	}
}

func WithDialTimeout(timeout time.Duration) OptionFunc {
	return func(c *Client) error {
		c.dialTimeout = timeout
		return nil
	}
}

func WithWriteTimeout(timeout time.Duration) OptionFunc {
	return func(c *Client) error {
		c.writeTimeout = timeout
		return nil
	}
}

func WithReadTimeout(timeout time.Duration) OptionFunc {
	return func(c *Client) error {
		c.readTimeout = timeout
		return nil
	}
}
