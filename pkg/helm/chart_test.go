package helm

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var testTarballPath = "../../testdata/charts/mychart/mychart-0.1.0.tgz"

func TestSetVersion(t *testing.T) {
	c, err := GetChartByName(testTarballPath)
	if err != nil {
		t.Error("unexpected error getting test tarball chart", err)
	}
	c.SetVersion("latest")
	if c.Metadata.Version != "latest" {
		t.Errorf("expected chart version to be latest, instead got %s", c.Metadata.Version)
	}
}

func TestGetChartByName(t *testing.T) {
	// Bad name
	_, err := GetChartByName("/non/existant/path/mychart-0.1.0.tgz")
	if err == nil {
		t.Error("expected error getting chart with bad name, instead got nil")
	}

	// Valid name
	c, err := GetChartByName(testTarballPath)
	if err != nil {
		t.Error("unexpected error getting test tarball chart", err)
	}
	if c.Metadata.Name != "mychart" {
		t.Errorf("expexted chart name to be mychart, instead got %s", c.Metadata.Name)
	}
	if c.Metadata.Version != "0.1.0" {
		t.Errorf("expexted chart version to be 0.1.0, instead got %s", c.Metadata.Version)
	}
}

func TestCreateChartPackage(t *testing.T) {
	c, err := GetChartByName(testTarballPath)
	if err != nil {
		t.Error("unexpected error getting test tarball chart", err)
	}

	tmp, err := ioutil.TempDir("", "helm-push-test")
	if err != nil {
		t.Error("unexpected error creating temp test dir", err)
	}
	defer os.RemoveAll(tmp)

	chartPackagePath, err := CreateChartPackage(c, tmp)
	if err != nil {
		t.Error("unexpected error creating chart package", err)
	}

	expectedPath := path.Join(tmp, "mychart-0.1.0.tgz")
	if chartPackagePath != expectedPath {
		t.Errorf("expected chart path to be %s, but was %s", expectedPath, chartPackagePath)
	}
}
