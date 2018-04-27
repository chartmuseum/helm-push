package chartmuseum

import (
	"testing"
	"encoding/base64"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
)

func TestDownloadFile(t *testing.T) {
	basicAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != basicAuthHeader {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(200)
			w.Write([]byte("hello world"))
		}
	}))
	defer ts.Close()

	cmClient := NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
	)

	resp, err := cmClient.DownloadFile("testfile")
	if err != nil {
		t.Error("error downloading testfile", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("error reading response body", err)
	}
	if s := string(b); s != "hello world" {
		t.Error("testfile contents were incorrect", s)
	}

	// trigger request failure
	cmClient = NewClient(URL("kjebnrkvjbnerv"))
	_, err = cmClient.DownloadFile("testfile")
	if err == nil {
		t.Error("expecting error with bad auth instead got nil")
	}
}
