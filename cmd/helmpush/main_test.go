package main

import (
	"io/ioutil"
	"k8s.io/helm/pkg/getter"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/repo"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	settings        helm_env.EnvSettings
	testTarballPath = "../../testdata/charts/mychart/mychart-0.1.0.tgz"
)

func TestPushCmd(t *testing.T) {
	statusCode := 201
	body := "{\"success\": true}"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))
	defer ts.Close()

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
	entry.URL = ts.URL

	_, err = repo.NewChartRepository(&entry, getter.All(settings))
	if err != nil {
		t.Error("unexpected error created test repository", err)
	}

	f.Update(&entry)
	os.MkdirAll(home.Repository(), 0777)
	f.WriteFile(home.RepositoryFile(), 0644)

	os.Setenv("HELM_HOME", home.String())
	os.Setenv("HELM_REPO_USERNAME", "myuser")
	os.Setenv("HELM_REPO_PASSWORD", "mypass")
	os.Setenv("HELM_REPO_CONTEXT_PATH", "/x/y/z")

	// Not enough args
	args := []string{}
	cmd := newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err == nil {
		t.Error("expecting error with missing args, instead got nil")
	}

	// Bad chart path
	args = []string{"/this/this/not/a/chart", "helm-push-test"}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err == nil {
		t.Error("expecting error with bad chart path, instead got nil")
	}

	// Bad repo name
	args = []string{testTarballPath, "wkerjbnkwejrnkj"}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err == nil {
		t.Error("expecting error with bad repo name, instead got nil")
	}

	// Happy path
	args = []string{testTarballPath, "helm-push-test"}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err != nil {
		t.Error("unexpecting error uploading tarball", err)
	}

	// Trigger 409
	statusCode = 409
	body = "{\"error\": \"package already exists\"}"
	args = []string{testTarballPath, "helm-push-test"}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err == nil {
		t.Error("expecting error with 409, instead got nil")
	}

	// Unable to parse JSON response body
	statusCode = 500
	body = "qkewjrnvqejrnbvjern"
	args = []string{testTarballPath, "helm-push-test"}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err == nil {
		t.Error("expecting error with bad response body, instead got nil")
	}

	// cm:// downloader
	os.Setenv("HELM_REPO_USE_HTTP", "true")
	downloaderBaseURL := strings.Replace(ts.URL, "http://", "cm://", 1)

	// fails with no file path
	args = []string{"", "", "", downloaderBaseURL}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err == nil {
		t.Error("expecting error with bad cm:// url, instead got nil")
	}

	// index.yaml
	args = []string{"", "", "", downloaderBaseURL+"/index.yaml"}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err != nil {
		t.Error("unexpecting error trying to download index.yaml", err)
	}

	// charts/mychart-0.1.0.tgz
	args = []string{"", "", "", downloaderBaseURL+"/charts/mychart-0.1.0.tgz"}
	cmd = newPushCmd(args)
	err = cmd.RunE(cmd, args)
	if err != nil {
		t.Error("unexpecting error trying to download charts/mychart-0.1.0.tgz", err)
	}
}
