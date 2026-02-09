package helm

import (
	"fmt"

	v3chart "helm.sh/helm/v3/pkg/chart"
	v3loader "helm.sh/helm/v3/pkg/chart/loader"
	v3chartutil "helm.sh/helm/v3/pkg/chartutil"

	v4chart "helm.sh/helm/v4/pkg/chart/v2"
	v4loader "helm.sh/helm/v4/pkg/chart/loader"
	v4chartutil "helm.sh/helm/v4/pkg/chart/v2/util"
)

type (
	// Chart is a helm package that contains metadata
	Chart struct {
		v3Chart  *v3chart.Chart
		v4Chart  *v4chart.Chart
		version  HelmMajorVersion
		Metadata *Metadata
	}
)

// Metadata returns a unified metadata view
type Metadata struct {
	Name       string
	Version    string
	AppVersion string
}

// GetMetadata returns the chart metadata
func (c *Chart) GetMetadata() *Metadata {
	if c.version == HelmMajorVersion3 {
		return &Metadata{
			Name:       c.v3Chart.Metadata.Name,
			Version:    c.v3Chart.Metadata.Version,
			AppVersion: c.v3Chart.Metadata.AppVersion,
		}
	}
	return &Metadata{
		Name:       c.v4Chart.Metadata.Name,
		Version:    c.v4Chart.Metadata.Version,
		AppVersion: c.v4Chart.Metadata.AppVersion,
	}
}

// SetVersion overrides the chart version
func (c *Chart) SetVersion(version string) {
	if c.version == HelmMajorVersion3 {
		c.v3Chart.Metadata.Version = version
	} else {
		c.v4Chart.Metadata.Version = version
	}
	// Update the exposed Metadata
	c.Metadata = c.GetMetadata()
}

// SetAppVersion overrides the app version
func (c *Chart) SetAppVersion(appVersion string) {
	if c.version == HelmMajorVersion3 {
		c.v3Chart.Metadata.AppVersion = appVersion
	} else {
		c.v4Chart.Metadata.AppVersion = appVersion
	}
	// Update the exposed Metadata
	c.Metadata = c.GetMetadata()
}

// GetChartByName returns a chart by "name", which can be
// either a directory or .tgz package
func GetChartByName(name string) (*Chart, error) {
	version := HelmMajorVersionCurrent()

	if version == HelmMajorVersion3 {
		cc, err := v3loader.Load(name)
		if err != nil {
			return nil, err
		}
		chart := &Chart{v3Chart: cc, version: version}
		chart.Metadata = chart.GetMetadata()
		return chart, nil
	} else {
		cc, err := v4loader.Load(name)
		if err != nil {
			return nil, err
		}
		v4c, ok := cc.(*v4chart.Chart)
		if !ok {
			return nil, fmt.Errorf("failed to load chart as v4 chart type")
		}
		chart := &Chart{v4Chart: v4c, version: version}
		chart.Metadata = chart.GetMetadata()
		return chart, nil
	}
}

// CreateChartPackage creates a new .tgz package in directory
func CreateChartPackage(c *Chart, outDir string) (string, error) {
	if c.version == HelmMajorVersion3 {
		return v3chartutil.Save(c.v3Chart, outDir)
	} else {
		return v4chartutil.Save(c.v4Chart, outDir)
	}
}
