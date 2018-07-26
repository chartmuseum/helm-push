package chartmuseum

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cmClient, err := NewClient(
		URL("http://localhost:8080"),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
		Timeout(60),
		CAFile("../../testdata/tls/ca.crt"),
		KeyFile("../../testdata/tls/test_key.key"),
		CertFile("../../testdata/tls/test_cert.crt"),
		InsecureSkipVerify(true),
	)

	if err != nil {
		t.Fatalf("expect creating a client instance but met error: %s", err)
	}

	if cmClient.opts.url != "http://localhost:8080" {
		t.Errorf("expected url to be http://localhost:8080, got %v", cmClient.opts.url)
	}

	if cmClient.opts.username != "user" {
		t.Errorf("expected username to be user, got %v", cmClient.opts.username)
	}

	if cmClient.opts.password != "pass" {
		t.Errorf("expected password to be pass, got %v", cmClient.opts.password)
	}

	if cmClient.opts.contextPath != "/my/context/path" {
		t.Errorf("expected context path to be /my/context/path, got %v", cmClient.opts.contextPath)
	}

	if cmClient.opts.timeout != time.Minute {
		t.Errorf("expected timeout duration to be 1 minute, got %v", cmClient.opts.timeout)
	}

	if cmClient.opts.caFile != "../../testdata/tls/ca.crt" {
		t.Errorf("expected ca file path to be '../../testdata/tls/ca.crt' but got %v", cmClient.opts.caFile)
	}

	if cmClient.opts.certFile != "../../testdata/tls/test_cert.crt" {
		t.Errorf("expected cert file path to be '../../testdata/tls/test_cert.crt' but got %v", cmClient.opts.certFile)
	}

	if cmClient.opts.keyFile != "../../testdata/tls/test_key.key" {
		t.Errorf("expected key file path to be '../../testdata/tls/test_key.key' but got %v", cmClient.opts.keyFile)
	}

	if !cmClient.opts.InsecureSkipVerify {
		t.Errorf("expected insecure flag to be 'true' but got %v", cmClient.opts.InsecureSkipVerify)
	}
}
