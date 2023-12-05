package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	testApp := application{
		errLog:  log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
	}
	return &testApp
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// get method
func (ts *testServer) get(t *testing.T, urlPath string) (statusCode int, header http.Header, body string) {
	res, err := ts.Client().Get(ts.URL + "/" + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	rBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(rBody)
	return res.StatusCode, res.Header, string(rBody)
}
