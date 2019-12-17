package helm

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	v2chartutil "k8s.io/helm/pkg/chartutil"
	v2chart "k8s.io/helm/pkg/proto/hapi/chart"
)

type (
	// Chart is a helm package that contains metadata
	Chart struct {
		V3 *chart.Chart
		V2 *v2chart.Chart
	}
)

// SetVersion overrides the chart version
func (c *Chart) SetVersion(version string) {
	if c.V2 != nil {
		c.V2.Metadata.Version = version
	} else {
		c.V3.Metadata.Version = version
	}
}

// GetChartByName returns a chart by "name", which can be
// either a directory or .tgz package
func GetChartByName(name string) (*Chart, error) {
	c := &Chart{}
	v3c, err := loader.Load(name)
	if err != nil {
		return nil, err
	}

	// If the Helm 2 API version (v1) is detected, use the old
	// method to load the chart
	if v3c.Metadata.APIVersion == chart.APIVersionV1 {
		v2c, err := v2chartutil.Load(name)
		if err != nil {
			return nil, err
		}
		c.V2 = v2c
	} else {
		c.V3 = v3c
	}

	return c, nil
}

// CreateChartPackage creates a new .tgz package in directory
func CreateChartPackage(c *Chart, outDir string) (string, error) {
	if c.V2 != nil {
		return v2chartutil.Save(c.V2, outDir)
	} else {
		return chartutil.Save(c.V3, outDir)
	}
}
