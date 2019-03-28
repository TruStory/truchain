package truapi

import (
	"net/http"
	"os"
	"strings"

	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/twitter"
)

// IssueSession creates a session and redirects the logged in user to the correct page
func IssueSession(ta *TruAPI) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		addr, err := CalibrateUser(ta, twitterUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		twitterProfile := &db.TwitterProfile{
			ID:        twitterUser.ID,
			Address:   addr,
			Username:  twitterUser.ScreenName,
			FullName:  twitterUser.Name,
			AvatarURI: strings.Replace(twitterUser.ProfileImageURL, "_normal", "_bigger", 1),
		}
		// upserting the twitter profile
		err = ta.DBClient.UpsertTwitterProfile(twitterProfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		cookie, err := cookies.GetLoginCookie(twitterProfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, req, os.Getenv("AUTH_LOGIN_REDIR"), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// HandleOAuthFailure handles the failed oAuth requests gracefully
func HandleOAuthFailure(ta *TruAPI) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		// if the authorization was purposefully denied by the user
		if req.FormValue("denied") != "" {
			http.Redirect(w, req, os.Getenv("AUTH_DENIED_REDIR"), http.StatusFound)
			return
		}

		// if any other error
		ctx := req.Context()
		err := gologin.ErrorFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// should be unreachable, ErrorFromContext always returns some non-nil error
		http.Error(w, "", http.StatusInternalServerError)
	}
	return http.HandlerFunc(fn)
}
