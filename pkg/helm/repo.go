package helm

import (
	"fmt"
	urllib "net/url"
	"os"
	"strings"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
)

type (
	// Repo represents a collection of parameters for chart repository
	Repo struct {
		*repo.ChartRepository
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
	cr, err := repo.NewChartRepository(entry, getter.All(cli.New()))
	if err != nil {
		return nil, err
	}
	return &Repo{cr}, nil
}

// TempRepoFromURL builds a temporary Repo from a given URL
func TempRepoFromURL(url string) (*Repo, error) {
	u, err := urllib.Parse(url)
	if err != nil {
		return nil, err
	}
	entry := &repo.Entry{}
	if u.User != nil {
		// remove the username/password section from URL
		pass, _ := u.User.Password()
		entry.URL = strings.Split(url, "://")[0] + "://" + strings.Split(url, fmt.Sprintf("%s@", pass))[1]
		entry.Username = u.User.Username()
		entry.Password = pass
	} else {
		entry.URL = url
	}
	cr, err := repo.NewChartRepository(entry, getter.All(cli.New()))
	if err != nil {
		return nil, err
	}
	return &Repo{cr}, nil
}

func repoFile() (*repo.File, error) {
	var repoFilePath string
	helmMajorVersion := GetHelmMajorVersion()
	if helmMajorVersion == HelmMajorVersion2 {
		home := helmHome()
		repoFilePath = home.RepositoryFile()
	} else {
		// TODO: fix, this only works for default mac
		userHome, _ := os.UserHomeDir()
		repoFilePath = fmt.Sprintf("%s/Library/Preferences/helm/repositories.yaml", userHome)
	}
	repoFile, err := repo.LoadFile(repoFilePath)
	return repoFile, err
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

func findRepoEntry(name string, r *repo.File) (*repo.Entry, bool) {
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
