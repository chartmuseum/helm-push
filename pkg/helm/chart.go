package helm

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

type (
	// Chart is a helm package that contains metadata
	Chart struct {
		*chart.Chart
	}
)

// SetVersion overrides the chart version
func (c *Chart) SetVersion(version string) {
	c.Metadata.Version = version
}

// GetChartByName returns a chart by "name", which can be
// either a directory or .tgz package
func GetChartByName(name string) (*Chart, error) {
	cc, err := loader.Load(name)
	if err != nil {
		return nil, err
	}
	return &Chart{cc}, nil
}

// CreateChartPackage creates a new .tgz package in directory
func CreateChartPackage(c *Chart, outDir string) (string, error) {
	return chartutil.Save(c.Chart, outDir)
}
