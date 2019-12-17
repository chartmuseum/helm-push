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
	if HelmMajorVersionCurrent() == HelmMajorVersion2 {
		cc, err := v2chartutil.Load(name)
		if err != nil {
			return nil, err
		}
		c.V2 = cc
	} else {
		cc, err := loader.Load(name)
		if err != nil {
			return nil, err
		}
		c.V3 = cc
	}
	return c, nil
}

// CreateChartPackage creates a new .tgz package in directory
func CreateChartPackage(c *Chart, outDir string) (string, error) {
	if HelmMajorVersionCurrent() == HelmMajorVersion2 {
		return v2chartutil.Save(c.V2, outDir)
	} else {
		return chartutil.Save(c.V3, outDir)
	}
}
