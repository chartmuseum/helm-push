package chartmuseum

import (
	"fmt"
	"net/http"

	v2tlsutil "k8s.io/helm/pkg/tlsutil"
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
func NewClient(opts ...Option) (*Client, error) {
	var client Client
	client.Client = &http.Client{}
	client.Option(Timeout(30))
	client.Option(opts...)
	client.Timeout = client.opts.timeout

	//Enable tls config if configured
	tr, err := newTransport(
		client.opts.certFile,
		client.opts.keyFile,
		client.opts.caFile,
		client.opts.insecureSkipVerify,
	)
	if err != nil {
		return nil, err
	}

	client.Transport = tr

	return &client, nil
}

//Create transport with TLS config
func newTransport(certFile, keyFile, caFile string, insecureSkipVerify bool) (*http.Transport, error) {
	transport := &http.Transport{}

	tlsConf, err := v2tlsutil.NewClientTLS(certFile, keyFile, caFile)
	if err != nil {
		return nil, fmt.Errorf("can't create TLS config: %s", err.Error())
	}
	tlsConf.InsecureSkipVerify = insecureSkipVerify

	transport.TLSClientConfig = tlsConf
	transport.Proxy = http.ProxyFromEnvironment

	return transport, nil
}
