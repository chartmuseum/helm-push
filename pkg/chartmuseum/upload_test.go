package chartmuseum

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	testTarballPath = "../../testdata/charts/mychart/mychart-0.1.0.tgz"
)

func TestUploadChartPackage(t *testing.T) {
	chartUploaded := false

	basicAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.String(), "/my/context/path") {
			w.WriteHeader(404)
		} else if r.Header.Get("Authorization") != basicAuthHeader {
			w.WriteHeader(401)
		} else if chartUploaded {
			if _, ok := r.URL.Query()["force"]; ok {
				w.WriteHeader(201)
			} else {
				w.WriteHeader(409)
			}
		} else {
			chartUploaded = true
			w.WriteHeader(201)
		}
	}))
	defer ts.Close()

	// Happy path
	cmClient := NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
	)
	resp, err := cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("error uploading chart package", err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("expecting 201 instead got %d", resp.StatusCode)
	}

	// Attempt to re-upload without force, trigger 409
	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("error uploading chart package", err)
	}
	if resp.StatusCode != 409 {
		t.Errorf("expecting 409 instead got %d", resp.StatusCode)
	}

	// Upload with force
	resp, err = cmClient.UploadChartPackage(testTarballPath, true)
	if err != nil {
		t.Error("error uploading chart package", err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("expecting 201 instead got %d", resp.StatusCode)
	}

	// Bad package path
	resp, err = cmClient.UploadChartPackage("/non/existant/path/mychart-0.1.0.tgz", false)
	if err == nil {
		t.Error("expecting error with bad package path, instead got nil")
	}

	// Bad URL
	cmClient = NewClient(URL("jaswehfgew"))
	_, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err == nil {
		t.Error("expecting error with bad package path, instead got nil")
	}

	// Bad context path
	cmClient = NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/crappy/context/path"),
		Timeout(5),
	)
	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("unexpected error with bad context path", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("expecting 404 instead got %d", resp.StatusCode)
	}

	// Unauthorized, invalid user/pass combo (basic auth)
	cmClient = NewClient(
		URL(ts.URL),
		Username("baduser"),
		Password("badpass"),
		ContextPath("/my/context/path"),
	)
	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("unexpected error with invalid user/pass combo (basic auth)", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expecting 401 instead got %d", resp.StatusCode)
	}

	// Unauthorized, missing user/pass combo (basic auth)
	cmClient = NewClient(
		URL(ts.URL),
		ContextPath("/my/context/path"),
	)
	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("unexpected error with missing user/pass combo (basic auth)", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expecting 401 instead got %d", resp.StatusCode)
	}
}
