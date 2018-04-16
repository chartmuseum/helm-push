package chartmuseum

import (
	"net/http"
)

type (
	// Client is an HTTP client to connect to ChartMuseum
	Client struct {
		*http.Client
		opts options
	}
)

// Option configures the client with the provided options.
func (client *Client) Option(opts ...Option) *Client {
	for _, opt := range opts {
		opt(&client.opts)
	}
	return client
}

// NewClient creates a new client.
func NewClient(opts ...Option) *Client {
	var client Client
	client.Client = &http.Client{}
	client.Option(Timeout(30))
	client.Option(opts...)
	client.Timeout = client.opts.timeout
	return &client
}
