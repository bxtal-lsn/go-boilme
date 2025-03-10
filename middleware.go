package boilme

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/justinas/nosurf"
)

func (b *Boilme) SessionLoad(next http.Handler) http.Handler {
	return b.Session.LoadAndSave(next)
}

func (b *Boilme) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(b.config.cookie.secure)

	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: secure,
		SameSite: http.SameSiteStrictMode,
		Domain: b.config.cookie.domain,
	})

	return csrfHandler
}

func (b *Boilme) CheckForMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maintenanceMode {
			if !strings.Contains(r.URL.Path, "/public/maintenance.html") {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Header().Set("Retry-After:", "300")
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				http.ServeFile(w, r, fmt.Sprintf("%s/public/maintenance.html", b.RootPath))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}