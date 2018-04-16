package chartmuseum

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cmClient := NewClient(
		URL("http://localhost:8080"),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
		Timeout(60),
	)

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
}
