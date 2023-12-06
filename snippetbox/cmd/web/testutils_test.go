package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/johnmerga/Mastering_Go/snippetbox/internal/models/mocks"
)

func newTestApplication(t *testing.T) *application {
	// Create an instance of the template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	// And a form decoder.
	formDecoder := form.NewDecoder()
	// And a session manager instance. Note that we use the same settings as
	// production, except that we *don't* set a Store for the session manager.
	// If no store is set, the SCS package will default to using a transient
	// in-memory store, which is ideal for testing purposes.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	return &application{
		errLog:   log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
		snippets: &mocks.SnippetModel{},
		users:    &mocks.UserModel{},

		// Use the mock.
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	// Initialize a new cookie jar.
	// jar, err := cookiejar.New(nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// Add the cookie jar to the test server client. Any response cookies will
	// now be stored and sent with subsequent requests when using this client.
	// ts.Client().Jar = jar
	// ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
	// 	return http.ErrUseLastResponse
	// }
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

//	func newTestApplication(t *testing.T) *application {
//		testApp := application{
//			errLog:  log.New(io.Discard, "", 0),
//			infoLog: log.New(io.Discard, "", 0),
//		}
//		return &testApp
//	}