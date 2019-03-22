package truapi

import (
	"net/http"
	"os"

	"github.com/TruStory/truchain/x/truapi/cookies"
)

// Logout deletes a session and redirects the logged in user to the correct page
func Logout() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		cookie := cookies.GetLogoutCookie()
		http.SetCookie(w, cookie)
		http.Redirect(w, req, os.Getenv("OAUTH_SUCCESS_REDIR"), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
