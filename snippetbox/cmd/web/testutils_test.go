package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
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
	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the cookie jar to the test server client. Any response cookies will
	// now be stored and sent with subsequent requests when using this client.
	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx
	// response is received by the client, and by always returning a
	// http.ErrUseLastResponse error it forces the client to immediately return
	// the received response.
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
