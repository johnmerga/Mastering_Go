package main

import (
	"net/http"
	"testing"

	"github.com/johnmerga/Mastering_Go/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	statusCode, _, body := ts.get(t, "/ping")
	assert.Equals(t, statusCode, http.StatusOK)
	assert.Equals(t, body, "OK")
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
