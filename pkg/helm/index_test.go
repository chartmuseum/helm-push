package helm

import (
	"testing"
)

func TestLoadIndex(t *testing.T) {
	// No context path
	index, err := LoadIndex([]byte("apiVersion: v1\nentries: {}\ngenerated: \"2018-08-08T08:21:33Z\"\n"))
	if err != nil {
		t.Error("unexpected error loading index", err)
	}
	if index.ContextPath != "" {
		t.Errorf("expexted empty context path, instead got %s", index.ContextPath)
	}

	// Has context path
	index, err = LoadIndex([]byte("apiVersion: v1\ncontextPath: /helm/v1\nentries: {}\ngenerated: \"2018-08-08T08:21:33Z\"\n"))
	if err != nil {
		t.Error("unexpected error loading index", err)
	}
	if index.ContextPath != "/helm/v1" {
		t.Errorf("expexted context path to be /helm/v1, instead got %s", index.ContextPath)
	}
}
