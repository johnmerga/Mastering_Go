package main

	"github.com/fatih/color"

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methodColor := color.New(color.FgHiMagenta).Add(color.Bold)
		app.infoLog.Printf("- %s - %s - %s - %s", r.RemoteAddr, r.Proto, methodColor.Sprint(r.Method), r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
