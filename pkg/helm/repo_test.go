package helm

import (
	"io/ioutil"
	"k8s.io/helm/pkg/getter"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/repo"
	"os"
	"testing"
)

var (
	settings helm_env.EnvSettings
)

func TestGetRepoByName(t *testing.T) {
	// Non-existant repo
	_, err := GetRepoByName("nonexistantrepo")
	if err == nil {
		t.Error("expecting error with bad repo name, instead got nil")
	}

	// Create new Helm home w/ test repo
	tmp, err := ioutil.TempDir("", "helm-push-test")
	if err != nil {
		t.Error("unexpected error creating temp test dir", err)
	}
	defer os.RemoveAll(tmp)

	home := helmpath.Home(tmp)
	f := repo.NewRepoFile()

	entry := repo.Entry{}
	entry.Name = "helm-push-test"
	entry.URL = "http://localhost:8080"

	_, err = repo.NewChartRepository(&entry, getter.All(settings))
	if err != nil {
		t.Error("unexpected error created test repository", err)
	}

	f.Update(&entry)
	os.MkdirAll(home.Repository(), 0777)
	f.WriteFile(home.RepositoryFile(), 0644)

	os.Setenv("HELM_HOME", home.String())

	// Retrieve test repo
	_, err = GetRepoByName("helm-push-test")
	if err != nil {
		t.Error("unexpected error getting test repo", err)
	}

	// Err, missing repofile
	os.RemoveAll(tmp)
	_, err = GetRepoByName("helm-push-test")
	if err == nil {
		t.Error("expecting error getting test repo after removed, instead got nil")
	}

}

func TestTempRepoFromURL(t *testing.T) {
	url := "https://my.chart.repo.com"
	repo, err := TempRepoFromURL(url)
	if err != nil {
		t.Error("unexpected error getting temp repo from URL", err)
	}
	if repo.URL != url {
		t.Error("expecting repo URL to match what was provided")
	}

	url = "https://user:p@ss@my.chart.repo.com/a/b/c/"
	repo, err = TempRepoFromURL(url)
	if err != nil {
		t.Error("unexpected error getting temp repo from URL, with basic auth", err)
	}
	if repo.URL != "https://my.chart.repo.com/a/b/c/" {
		t.Error("expecting repo URL to have basic auth removed")
	}
	if repo.Username != "user" {
		t.Error("expecting repo username to be extracted from URL")
	}
	if repo.Password != "p@ss" {
		t.Error("expecting repo password to be extracted from URL")
	}
}
