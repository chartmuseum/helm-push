package helm

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"k8s.io/helm/pkg/repo"
)

type (
	// Index represents the index file in a chart repository
	Index struct {
		*repo.IndexFile
		ContextPath string `json:"contextPath"`
	}

	// IndexDownloader is a function to download the index
	IndexDownloader func() ([]byte, error)
)

// GetIndexByRepo returns index by repository
func GetIndexByRepo(repo *Repo, downloadIndex IndexDownloader) (*Index, error) {
	if repo.Cache != "" {
		return GetIndexByDownloader(func() ([]byte, error) {
			return ioutil.ReadFile(repo.Cache)
		})
	}
	return GetIndexByDownloader(downloadIndex)
}

// GetIndexByDownloader takes binary data from IndexDownloader and returns an Index object
func GetIndexByDownloader(downloadIndex IndexDownloader) (*Index, error) {
	b, err := downloadIndex()
	if err != nil {
		return nil, err
	}
	return LoadIndex(b)
}

// LoadIndex loads an index file
func LoadIndex(data []byte) (*Index, error) {
	i := &Index{}
	if err := yaml.Unmarshal(data, i); err != nil {
		return i, err
	}
	i.SortEntries()
	return i, nil
}
