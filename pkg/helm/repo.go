package helm

import (
	"fmt"
	urllib "net/url"
	"strings"

	v3cli "helm.sh/helm/v3/pkg/cli"
	v3getter "helm.sh/helm/v3/pkg/getter"
	v3repo "helm.sh/helm/v3/pkg/repo"

	v4cli "helm.sh/helm/v4/pkg/cli"
	v4getter "helm.sh/helm/v4/pkg/getter"
	v4repo "helm.sh/helm/v4/pkg/repo/v1"
)

type (
	// Repo represents a collection of parameters for chart repository
	Repo struct {
		v3Repo  *v3repo.ChartRepository
		v4Repo  *v4repo.ChartRepository
		version HelmMajorVersion
	}
)

// Accessor methods for Config fields
func (r *Repo) GetConfigName() string {
	if r.version == HelmMajorVersion3 {
		return r.v3Repo.Config.Name
	}
	return r.v4Repo.Config.Name
}

func (r *Repo) GetConfigURL() string {
	if r.version == HelmMajorVersion3 {
		return r.v3Repo.Config.URL
	}
	return r.v4Repo.Config.URL
}

func (r *Repo) GetConfigUsername() string {
	if r.version == HelmMajorVersion3 {
		return r.v3Repo.Config.Username
	}
	return r.v4Repo.Config.Username
}

func (r *Repo) GetConfigPassword() string {
	if r.version == HelmMajorVersion3 {
		return r.v3Repo.Config.Password
	}
	return r.v4Repo.Config.Password
}

func (r *Repo) GetCachePath() string {
	if r.version == HelmMajorVersion3 {
		return r.v3Repo.CachePath
	}
	return r.v4Repo.CachePath
}

// GetRepoByName returns repository by name
func GetRepoByName(name string) (*Repo, error) {
	version := HelmMajorVersionCurrent()

	if version == HelmMajorVersion3 {
		r, err := repoFileV3()
		if err != nil {
			return nil, err
		}
		entry, exists := findRepoEntryV3(name, r)
		if !exists {
			return nil, fmt.Errorf("no repo named %q found", name)
		}

		settings := v3cli.New()
		getters := v3getter.All(settings)
		cr, err := v3repo.NewChartRepository(entry, getters)
		if err != nil {
			return nil, err
		}

		return &Repo{v3Repo: cr, version: version}, nil
	} else {
		r, err := repoFileV4()
		if err != nil {
			return nil, err
		}
		entry, exists := findRepoEntryV4(name, r)
		if !exists {
			return nil, fmt.Errorf("no repo named %q found", name)
		}

		settings := v4cli.New()
		getters := v4getter.All(settings)
		cr, err := v4repo.NewChartRepository(entry, getters)
		if err != nil {
			return nil, err
		}

		return &Repo{v4Repo: cr, version: version}, nil
	}
}

// TempRepoFromURL builds a temporary Repo from a given URL
func TempRepoFromURL(url string) (*Repo, error) {
	version := HelmMajorVersionCurrent()
	u, err := urllib.Parse(url)
	if err != nil {
		return nil, err
	}

	if version == HelmMajorVersion3 {
		entry := &v3repo.Entry{}
		if u.User != nil {
			// remove the username/password section from URL
			pass, _ := u.User.Password()
			entry.URL = strings.Split(url, "://")[0] + "://" + strings.Split(url, fmt.Sprintf("%s@", pass))[1]
			entry.Username = u.User.Username()
			entry.Password = pass
		} else {
			entry.URL = url
		}
		cr, err := v3repo.NewChartRepository(entry, v3getter.All(v3cli.New()))
		if err != nil {
			return nil, err
		}
		return &Repo{v3Repo: cr, version: version}, nil
	} else {
		entry := &v4repo.Entry{}
		if u.User != nil {
			// remove the username/password section from URL
			pass, _ := u.User.Password()
			entry.URL = strings.Split(url, "://")[0] + "://" + strings.Split(url, fmt.Sprintf("%s@", pass))[1]
			entry.Username = u.User.Username()
			entry.Password = pass
		} else {
			entry.URL = url
		}
		cr, err := v4repo.NewChartRepository(entry, v4getter.All(v4cli.New()))
		if err != nil {
			return nil, err
		}
		return &Repo{v4Repo: cr, version: version}, nil
	}
}

func repoFileV3() (*v3repo.File, error) {
	settings := v3cli.New()
	repoFilePath := settings.RepositoryConfig
	return v3repo.LoadFile(repoFilePath)
}

func repoFileV4() (*v4repo.File, error) {
	settings := v4cli.New()
	repoFilePath := settings.RepositoryConfig
	return v4repo.LoadFile(repoFilePath)
}

func findRepoEntryV3(name string, r *v3repo.File) (*v3repo.Entry, bool) {
	for _, re := range r.Repositories {
		if re.Name == name {
			return re, true
		}
	}
	return nil, false
}

func findRepoEntryV4(name string, r *v4repo.File) (*v4repo.Entry, bool) {
	for _, re := range r.Repositories {
		if re.Name == name {
			return re, true
		}
	}
	return nil, false
}
