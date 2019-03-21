package truapi

import (
	"net/http"
	"os"
	"time"
)

// Logout deletes a session and redirects the logged in user to the correct page
func Logout() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		cookie := http.Cookie{
			Name:     "tru-user",
			HttpOnly: true,
			Value:    "",
			Expires:  time.Now(),
			Domain:   os.Getenv("COOKIE_HOST"),
			MaxAge:   0,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, req, os.Getenv("OAUTH_SUCCESS_REDIR"), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
