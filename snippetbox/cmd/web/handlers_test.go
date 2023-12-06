package main

import (
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/johnmerga/Mastering_Go/snippetbox/internal/assert"
	"github.com/johnmerga/Mastering_Go/snippetbox/internal/models/mocks"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	statusCode, _, body := ts.get(t, "/ping")
	assert.Equals(t, statusCode, http.StatusOK)
	assert.Equals(t, body, "OK")
}
func newTestApplication1(t *testing.T) *application {
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
func TestSnippetView(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{name: "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equals(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

// func TestPing(t *testing.T) {
// 	rr := httptest.NewRecorder()
// 	r, err := http.NewRequest(http.MethodGet, "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ping(rr, r)
// 	rs := rr.Result()
// 	assert.Equals(t, rs.StatusCode, http.StatusOK)
// 	defer rs.Body.Close()
// 	body, err := io.ReadAll(rs.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bytes.TrimSpace(body)
// 	assert.Equals(t, string(body), "OK")
// }

// func TestPing(t *testing.T) {
// 	app := &application{
// 		// The reason for this is that the loggers are
// 		// needed by the logRequest and recoverPanic middlewares, which are used by our
// 		// application on every route.
// 		errLog:  log.New(io.Discard, "", 0),
// 		infoLog: log.New(io.Discard, "", 0),
// 	}
//
// 	testServer := httptest.NewTLSServer(app.routes())
// 	defer testServer.Close()
// 	res, err := testServer.Client().Get(testServer.URL + "/ping")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equals(t, res.StatusCode, http.StatusOK)
// 	defer res.Body.Close()
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bytes.TrimSpace(body)
// 	assert.Equals(t, string(body), "OK")
//