package chartmuseum

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
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

	tlsConf, err := newClientTLS(certFile, keyFile, caFile)
	if err != nil {
		return nil, fmt.Errorf("can't create TLS config: %s", err.Error())
	}
	tlsConf.InsecureSkipVerify = insecureSkipVerify

	transport.TLSClientConfig = tlsConf
	transport.Proxy = http.ProxyFromEnvironment

	return transport, nil
}

// newClientTLS creates a TLS configuration for HTTP client
func newClientTLS(certFile, keyFile, caFile string) (*tls.Config, error) {
	config := &tls.Config{}

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("can't load key pair from cert %s and key %s: %w", certFile, keyFile, err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("can't read CA cert file %s: %w", caFile, err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA cert from %s", caFile)
		}
		config.RootCAs = caCertPool
	}

	return config, nil
}
