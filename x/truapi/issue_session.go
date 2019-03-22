package truapi

import (
	"net/http"
	"os"

	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"
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
			AvatarURI: twitterUser.ProfileImageURL,
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
		http.Redirect(w, req, os.Getenv("OAUTH_SUCCESS_REDIR"), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
