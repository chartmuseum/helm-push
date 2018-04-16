package helm

import (
	"fmt"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/repo"
	"os"
)

type (
	// Repo represents a collection of parameters for chart repository
	Repo struct {
		*repo.Entry
	}
)

// GetRepoByName returns repository by name
func GetRepoByName(name string) (*Repo, error) {
	r, err := repoFile()
	if err != nil {
		return nil, err
	}
	entry, exists := findRepoEntry(name, r)
	if !exists {
		return nil, fmt.Errorf("no repo named %q found", name)
	}
	return &Repo{entry}, nil
}

func repoFile() (*repo.RepoFile, error) {
	home := helmHome()
	return repo.LoadRepositoriesFile(home.RepositoryFile())
}

func helmHome() helmpath.Home {
	var helmHomePath string
	if v, ok := os.LookupEnv("HELM_HOME"); ok {
		helmHomePath = v
	} else {
		helmHomePath = helm_env.DefaultHelmHome
	}
	return helmpath.Home(helmHomePath)
}

func findRepoEntry(name string, r *repo.RepoFile) (*repo.Entry, bool) {
	var entry *repo.Entry
	exists := false
	for _, re := range r.Repositories {
		if re.Name == name {
			entry = re
			exists = true
			break
		}
	}
	return entry, exists
}
