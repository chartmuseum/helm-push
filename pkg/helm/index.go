package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	v3repo "helm.sh/helm/v3/pkg/repo"
	v4repo "helm.sh/helm/v4/pkg/repo/v1"
)

type (
	// Index represents the index file in a chart repository
	Index struct {
		v3Index    *v3repo.IndexFile
		v4Index    *v4repo.IndexFile
		version    HelmMajorVersion
		ServerInfo ServerInfo `json:"serverInfo"`
	}

	// IndexDownloader is a function to download the index
	IndexDownloader func() ([]byte, error)
)

// SortEntries sorts entries in the index
func (i *Index) SortEntries() {
	if i.version == HelmMajorVersion3 {
		i.v3Index.SortEntries()
	} else {
		i.v4Index.SortEntries()
	}
}

// GetIndexByRepo returns index by repository
func GetIndexByRepo(repo *Repo, downloadIndex IndexDownloader) (*Index, error) {
	configName := repo.GetConfigName()
	cachePath := repo.GetCachePath()

	if configName != "" {
		return GetIndexByDownloader(func() ([]byte, error) {
			return os.ReadFile(filepath.Join(cachePath, fmt.Sprintf("%s-index.yaml", configName)))
		}, repo.version)
	}
	return GetIndexByDownloader(downloadIndex, repo.version)
}

// GetIndexByDownloader takes binary data from IndexDownloader and returns an Index object
func GetIndexByDownloader(downloadIndex IndexDownloader, version HelmMajorVersion) (*Index, error) {
	b, err := downloadIndex()
	if err != nil {
		return nil, err
	}
	return LoadIndex(b, version)
}

// LoadIndex loads an index file
func LoadIndex(data []byte, version HelmMajorVersion) (*Index, error) {
	i := &Index{version: version}

	if version == HelmMajorVersion3 {
		i.v3Index = &v3repo.IndexFile{}
		if err := yaml.Unmarshal(data, i.v3Index); err != nil {
			return i, err
		}
		// Also unmarshal ServerInfo
		if err := yaml.Unmarshal(data, i); err != nil {
			return i, err
		}
		i.v3Index.SortEntries()
	} else {
		i.v4Index = &v4repo.IndexFile{}
		if err := yaml.Unmarshal(data, i.v4Index); err != nil {
			return i, err
		}
		// Also unmarshal ServerInfo
		if err := yaml.Unmarshal(data, i); err != nil {
			return i, err
		}
		i.v4Index.SortEntries()
	}

	return i, nil
}
