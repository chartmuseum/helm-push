package chartmuseum

import (
	"crypto/tls"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	basicAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != basicAuthHeader {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(200)
			_, err := w.Write([]byte("hello world"))
			if err != nil {
				return
			}

		}

	}))
	defer ts.Close()

	cmClient, err := NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
	)
	if err != nil {
		t.Fatalf("expect creating a client instance but met error: %s", err)
	}

	resp, err := cmClient.DownloadFile("testfile")
	if err != nil {
		t.Fatal("error downloading testfile", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading response body", err)
	}
	if s := string(b); s != "hello world" {
		t.Fatal("testfile contents were incorrect", s)
	}

	// trigger request failure
	cmClient, err = NewClient(URL("kjebnrkvjbnerv"))
	if err != nil {
		t.Fatalf("expect creating a client instance but met error: %s", err)
	}
	_, err = cmClient.DownloadFile("testfile")
	if err == nil {
		t.Fatal("expecting error with bad auth instead got nil")
	}
}

func TestDownloadFileFromTlsServer(t *testing.T) {
	basicAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != basicAuthHeader {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(200)
			_, err := w.Write([]byte("hello world"))
			if err != nil {
				return
			}
		}
	}))
	cert, err := tls.LoadX509KeyPair(testServerCertPath, testServerKeyPath)
	if err != nil {
		t.Fatalf("failed to load certificate and key with error: %s", err.Error())
	}

	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	ts.StartTLS()
	defer ts.Close()

	//without ca file
	cmClient, err := NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
	)
	if err != nil {
		t.Fatalf("[without ca file] expect creating a client instance but met error: %s", err)
	}
	_, err = cmClient.DownloadFile("testfile")
	if err == nil {
		t.Error("[without ca file] expected error when downloading testfile without ca but got nil")
	}

	//with ca file
	cmClient, err = NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		CAFile(testServerCAPath),
	)
	if err != nil {
		t.Fatalf("[with ca file] expect creating a client instance but met error: %s", err)
	}

	resp, err := cmClient.DownloadFile("testfile")
	if err != nil {
		t.Fatalf("[with ca file] error downloading testfile: %s", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("[with ca file] error reading response body: %s", err)
	}
	if s := string(b); s != "hello world" {
		t.Fatalf("[with ca file] expected 'hello world' but got '%s'", s)
	}
}
