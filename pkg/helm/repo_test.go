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
