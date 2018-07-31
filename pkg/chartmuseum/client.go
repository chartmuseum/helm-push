package chartmuseum

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"k8s.io/helm/pkg/tlsutil"
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

	tlsConf, err := NewClientTLS(certFile, keyFile, caFile, insecureSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("can't create TLS config: %s", err.Error())
	}

	transport.TLSClientConfig = tlsConf
	transport.Proxy = http.ProxyFromEnvironment

	return transport, nil
}

//The fix for CA file has not been included in any formal release,
//so copy code here. Once it's released, we can use method 'NewClientTLS'
//in the pkg/tlsutil package to replace this copy code.
//For more details, please refer https://github.com/helm/helm/pull/3258
func newTLSConfigCommon(certFile, keyFile, caFile string, insecureSkipVerify bool) (*tls.Config, error) {
	config := tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}

	if certFile != "" && keyFile != "" {
		cert, err := tlsutil.CertFromFilePair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		config.Certificates = []tls.Certificate{*cert}
	}

	if !insecureSkipVerify && caFile != "" {
		cp, err := tlsutil.CertPoolFromFile(caFile)
		if err != nil {
			return nil, err
		}
		config.RootCAs = cp
	}

	return &config, nil
}

// NewClientTLS returns tls.Config appropriate for client auth.
func NewClientTLS(certFile, keyFile, caFile string, insecureSkipVerify bool) (*tls.Config, error) {
	return newTLSConfigCommon(certFile, keyFile, caFile, insecureSkipVerify)
}
