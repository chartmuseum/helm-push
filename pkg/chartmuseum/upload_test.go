package chartmuseum

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"k8s.io/helm/pkg/tlsutil"
)

var (
	testTarballPath    = "../../testdata/charts/helm2/mychart/mychart-0.1.0.tgz"
	testCertPath       = "../../testdata/tls/test_cert.crt"
	testKeyPath        = "../../testdata/tls/test_key.key"
	testCAPath         = "../../testdata/tls/ca.crt"
	testServerCAPath   = "../../testdata/tls/server_ca.crt"
	testServerCertPath = "../../testdata/tls/test_server.crt"
	testServerKeyPath  = "../../testdata/tls/test_server.key"
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
	cmClient, err := NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
	)
	if err != nil {
		t.Fatalf("[happy path] expect creating a client instance but met error: %s", err)
	}
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
	cmClient, _ = NewClient(URL("jaswehfgew"))
	_, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err == nil {
		t.Error("[bad URL] expecting error with bad package path, instead got nil")
	}

	// Bad context path
	cmClient, err = NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/crappy/context/path"),
		Timeout(5),
	)
	if err != nil {
		t.Fatalf("[bad context path] expect creating a client instance but met error: %s", err)
	}

	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("unexpected error with bad context path", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("expecting 404 instead got %d", resp.StatusCode)
	}

	// Unauthorized, invalid user/pass combo (basic auth)
	cmClient, err = NewClient(
		URL(ts.URL),
		Username("baduser"),
		Password("badpass"),
		ContextPath("/my/context/path"),
	)
	if err != nil {
		t.Fatalf("[unauthorized: invalid user/pass] expect creating a client instance but met error: %s", err)
	}
	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("unexpected error with invalid user/pass combo (basic auth)", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expecting 401 instead got %d", resp.StatusCode)
	}

	// Unauthorized, missing user/pass combo (basic auth)
	cmClient, err = NewClient(
		URL(ts.URL),
		ContextPath("/my/context/path"),
	)
	if err != nil {
		t.Fatalf("[unauthorized: missing user/pass] expect creating a client instance but met error: %s", err)
	}
	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Error("unexpected error with missing user/pass combo (basic auth)", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expecting 401 instead got %d", resp.StatusCode)
	}
}

func TestUploadChartPackageWithTlsServer(t *testing.T) {
	basicAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.String(), "/my/context/path") {
			w.WriteHeader(404)
		} else if r.Header.Get("Authorization") != basicAuthHeader {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(201)
		}
	}))

	cert, err := tls.LoadX509KeyPair(testCertPath, testKeyPath)
	if err != nil {
		t.Fatalf("failed to load certificate and key with error: %s", err.Error())
	}

	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	ts.StartTLS()
	defer ts.Close()

	//Without certificate
	cmClient, err := NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
	)
	if err != nil {
		t.Fatalf("[without certificate] expect creating a client instance but met error: %s", err)
	}

	_, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err == nil {
		t.Fatal("expected error returned when uploading package without cert to tls enabled https server")
	}

	//Enable insecure flag
	cmClient, err = NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
		InsecureSkipVerify(true),
	)
	if err != nil {
		t.Fatalf("[enable insecure flag] expect creating a client instance but met error: %s", err)
	}

	resp, err := cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Fatalf("[enable insecure flag] expected nil error but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("[enable insecure flag] expect status code 201 but got %d", resp.StatusCode)
	}

	//Upload with ca file
	cmClient, err = NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
		CAFile(testCAPath),
	)
	if err != nil {
		t.Fatalf("[upload with ca file] expect creating a client instance but met error: %s", err)
	}

	resp, err = cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Fatalf("[upload with ca file] expected nil error but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("[upload with ca file] expect status code 201 but got %d", resp.StatusCode)
	}
}

func TestUploadChartPackageWithVerifyingClientCert(t *testing.T) {
	basicAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.String(), "/my/context/path") {
			w.WriteHeader(404)
		} else if r.Header.Get("Authorization") != basicAuthHeader {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(201)
		}
	}))

	cert, err := tls.LoadX509KeyPair(testCertPath, testKeyPath)
	if err != nil {
		t.Fatalf("failed to load certificate and key with error: %s", err.Error())
	}

	caCertPool, err := tlsutil.CertPoolFromFile(testServerCAPath)
	if err != nil {
		t.Fatalf("load server CA file failed with error: %s", err.Error())
	}

	ts.TLS = &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		Rand:         rand.Reader,
	}
	ts.StartTLS()
	defer ts.Close()

	//Upload with cert and key files
	cmClient, err := NewClient(
		URL(ts.URL),
		Username("user"),
		Password("pass"),
		ContextPath("/my/context/path"),
		KeyFile(testServerKeyPath),
		CertFile(testServerCertPath),
		CAFile(testCAPath),
	)
	if err != nil {
		t.Fatalf("[upload with cert and key files] expect creating a client instance but met error: %s", err)
	}

	resp, err := cmClient.UploadChartPackage(testTarballPath, false)
	if err != nil {
		t.Fatalf("[upload with cert and key files] expected nil error but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("[upload with cert and key files] expect status code 201 but got %d", resp.StatusCode)
	}
}
