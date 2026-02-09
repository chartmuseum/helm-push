package helm

import (
	"bytes"
	"os"
	"testing"

	v3plugin "helm.sh/helm/v3/pkg/plugin"
	yamlv3 "go.yaml.in/yaml/v3"
	sigsyaml "sigs.k8s.io/yaml"
)

// helmV4MetadataV1 mirrors Helm v4's internal MetadataV1 struct.
// Helm v4 uses strict YAML parsing (KnownFields) for v1 manifests,
// so any field not in this struct will cause an unmarshal error.
type helmV4MetadataV1 struct {
	APIVersion    string         `yaml:"apiVersion"`
	Name          string         `yaml:"name"`
	Type          string         `yaml:"type"`
	Runtime       string         `yaml:"runtime"`
	Version       string         `yaml:"version"`
	SourceURL     string         `yaml:"sourceURL,omitempty"`
	Config        map[string]any `yaml:"config"`
	RuntimeConfig map[string]any `yaml:"runtimeConfig"`
}

// helmV4MetadataLegacy mirrors Helm v4's internal MetadataLegacy struct.
// Helm v4 uses non-strict YAML parsing for legacy manifests (no apiVersion).
type helmV4MetadataLegacy struct {
	Name            string                          `yaml:"name"`
	Version         string                          `yaml:"version"`
	Usage           string                          `yaml:"usage"`
	Description     string                          `yaml:"description"`
	PlatformCommand []helmV4PlatformCommand          `yaml:"platformCommand"`
	Command         string                          `yaml:"command"`
	IgnoreFlags     bool                            `yaml:"ignoreFlags"`
	PlatformHooks   map[string][]helmV4PlatformCommand `yaml:"platformHooks"`
	Hooks           map[string]string               `yaml:"hooks"`
	Downloaders     []helmV4Downloaders              `yaml:"downloaders"`
}

type helmV4PlatformCommand struct {
	OperatingSystem string   `yaml:"os"`
	Architecture    string   `yaml:"arch"`
	Command         string   `yaml:"command"`
	Args            []string `yaml:"args"`
}

type helmV4Downloaders struct {
	Protocols []string `yaml:"protocols"`
	Command   string   `yaml:"command"`
}

func TestPluginManifestHelm3Compatible(t *testing.T) {
	// Helm 3 uses sigs.k8s.io/yaml.UnmarshalStrict which converts YAML to JSON
	// and rejects unknown fields. This is the exact parsing path used by
	// helm.sh/helm/v3/pkg/plugin.LoadDir().
	files := map[string]string{
		"plugin.yaml":              "../../plugin.yaml",
		"testdata/plugin-helm3.yaml": "../../testdata/plugin-helm3.yaml",
	}

	for name, path := range files {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read %s: %v", path, err)
			}

			var m v3plugin.Metadata
			if err := sigsyaml.UnmarshalStrict(data, &m); err != nil {
				t.Errorf("%s is not valid for Helm 3: %v", name, err)
			}

			if m.Name != "cm-push" {
				t.Errorf("expected plugin name 'cm-push', got %q", m.Name)
			}
		})
	}
}

func TestPluginManifestHelm4V1Compatible(t *testing.T) {
	// Helm 4 uses go.yaml.in/yaml/v3 with KnownFields(true) for v1 manifests
	// (those with apiVersion set). This rejects any YAML field not defined
	// in MetadataV1. This is the exact parsing path used by
	// helm.sh/helm/v4/internal/plugin.loadMetadataV1().
	path := "../../testdata/plugin-helm4.yaml"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}

	var mv1 helmV4MetadataV1
	d := yamlv3.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	if err := d.Decode(&mv1); err != nil {
		t.Errorf("testdata/plugin-helm4.yaml is not valid for Helm 4 (v1 format): %v", err)
	}

	if mv1.Name != "cm-push" {
		t.Errorf("expected plugin name 'cm-push', got %q", mv1.Name)
	}
	if mv1.APIVersion != "v1" {
		t.Errorf("expected apiVersion 'v1', got %q", mv1.APIVersion)
	}
}

func TestPluginManifestHelm4LegacyCompatible(t *testing.T) {
	// Helm 4 supports legacy plugin manifests (no apiVersion field).
	// It uses non-strict YAML parsing for these. The shipped plugin.yaml
	// must be parseable by Helm 4's legacy loader since it's the first
	// file Helm sees before the install hook can swap manifests.
	// This is the exact parsing path used by
	// helm.sh/helm/v4/internal/plugin.loadMetadataLegacy().
	files := map[string]string{
		"plugin.yaml":              "../../plugin.yaml",
		"testdata/plugin-helm3.yaml": "../../testdata/plugin-helm3.yaml",
	}

	for name, path := range files {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read %s: %v", path, err)
			}

			// Verify no apiVersion field (which would route to strict V1 parsing)
			var peek struct {
				APIVersion string `yaml:"apiVersion"`
			}
			if err := yamlv3.Unmarshal(data, &peek); err != nil {
				t.Fatalf("failed to peek apiVersion: %v", err)
			}
			if peek.APIVersion != "" {
				t.Errorf("%s has apiVersion=%q; shipped plugin.yaml must use legacy format (no apiVersion) for Helm 3/4 compatibility", name, peek.APIVersion)
			}

			// Parse as legacy (non-strict, matching Helm 4's loadMetadataLegacy)
			var ml helmV4MetadataLegacy
			d := yamlv3.NewDecoder(bytes.NewReader(data))
			if err := d.Decode(&ml); err != nil {
				t.Errorf("%s is not valid for Helm 4 legacy format: %v", name, err)
			}

			if ml.Name != "cm-push" {
				t.Errorf("expected plugin name 'cm-push', got %q", ml.Name)
			}
		})
	}
}
