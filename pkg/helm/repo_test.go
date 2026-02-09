package helm

import (
	"io/ioutil"
	"os"
	"testing"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
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

	settings := cli.New()
	settings.RepositoryConfig = tmp + "/repositories.yaml"
	settings.RepositoryCache = tmp + "/repository"

	f := repo.NewFile()

	entry := &repo.Entry{}
	entry.Name = "helm-push-test"
	entry.URL = "http://localhost:8080"

	_, err = repo.NewChartRepository(entry, getter.All(settings))
	if err != nil {
		t.Error("unexpected error created test repository", err)
	}

	f.Update(entry)
	os.MkdirAll(tmp, 0777)
	f.WriteFile(settings.RepositoryConfig, 0644)

	os.Setenv("HELM_REPOSITORY_CONFIG", settings.RepositoryConfig)

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
	if repo.GetConfigURL() != url {
		t.Error("expecting repo URL to match what was provided")
	}

	url = "https://user:p@ss@my.chart.repo.com/a/b/c/"
	repo, err = TempRepoFromURL(url)
	if err != nil {
		t.Error("unexpected error getting temp repo from URL, with basic auth", err)
	}
	if repo.GetConfigURL() != "https://my.chart.repo.com/a/b/c/" {
		t.Error("expecting repo URL to have basic auth removed")
	}
	if repo.GetConfigUsername() != "user" {
		t.Error("expecting repo username to be extracted from URL")
	}
	if repo.GetConfigPassword() != "p@ss" {
		t.Error("expecting repo password to be extracted from URL")
	}
}
