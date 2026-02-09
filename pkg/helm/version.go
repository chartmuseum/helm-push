package helm

import (
	"os"
	"path/filepath"
	"strings"
)

type (
	HelmMajorVersion int
)

const (
	HelmMajorVersion3 = 3
	HelmMajorVersion4 = 4
)

var (
	helmMajorVersionCurrent HelmMajorVersion
)

func HelmMajorVersionCurrent() HelmMajorVersion {
	if helmMajorVersionCurrent != 0 {
		return helmMajorVersionCurrent
	}

	// Detect version by checking which plugin manifest format is being used
	// Helm 4 manifests contain "apiVersion", Helm 3 manifests don't
	if pluginDir := os.Getenv("HELM_PLUGIN_DIR"); pluginDir != "" {
		manifestFile := filepath.Join(pluginDir, "plugin.yaml")
		if data, err := os.ReadFile(manifestFile); err == nil {
			// Check if manifest contains "apiVersion" field (Helm 4 format)
			if strings.Contains(string(data), "apiVersion:") {
				helmMajorVersionCurrent = HelmMajorVersion4
				return helmMajorVersionCurrent
			}
		}
	}

	// Default to v3 if we can't determine from manifest
	helmMajorVersionCurrent = HelmMajorVersion3
	return helmMajorVersionCurrent
}
