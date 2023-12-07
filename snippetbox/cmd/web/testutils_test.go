package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/johnmerga/Mastering_Go/snippetbox/internal/models/mocks"
)

// Define a regular expression which captures the CSRF token value from the
// HTML for our user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type=\"hidden\" name=\"csrf_token\" value=\"(.+)\" .>`)

// var test := " <input type="hidden" name="csrf_token" value="(.+)'" />"
func extractCSRFToken(t *testing.T, body string) string {
	// Use the FindStringSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with the entire matched pattern in the
	// first position, and the values of any captured data in the subsequent
	// positions.
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	// Unescape the token value.
	return html.UnescapeString(string(matches[1]))
}

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
	res, err := ts.Client().Get(ts.URL + urlPath)
	println(ts.URL + "/" + urlPath)
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

func TestUserSignup(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running an end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)
	// Log the CSRF token value in our test output using the t.Logf() function.
	// The t.Logf() function works in the same way as fmt.Printf(), but writes
	// the provided message to the test output.
	t.Logf("CSRF token is: %q", csrfToken)
}

//	func newTestApplication(t *testing.T) *application {
//		testApp := application{
//			errLog:  log.New(io.Discard, "", 0),
//			infoLog: log.New(io.Discard, "", 0),
//		}
//		return &testApp
//	}
