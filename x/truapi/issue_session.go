package truapi

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/TruStory/truchain/x/db"
	"github.com/dghubble/gologin/twitter"
	"github.com/gorilla/securecookie"
	secp "github.com/tendermint/tendermint/crypto/secp256k1"
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

		// Fetch keypair of the user
		keyPair, err := ta.DBClient.KeyPairByTwitterProfileID(twitterUser.ID)
		if err != nil {
			panic(err)
		}

		// If not available, create new
		if keyPair.ID == 0 {
			privKey := secp.GenPrivKey()

			keyPair := &db.KeyPair{
				TwitterProfileID: twitterUser.ID,
				PrivateKey:       fmt.Sprintf("%x", privKey),
				PublicKey:        fmt.Sprintf("%x", privKey.PubKey().Bytes()),
			}
			err = ta.DBClient.Add(keyPair)
			if err != nil {
				panic(err)
			}
		}

		pubKeyBytes, _ := hex.DecodeString(keyPair.PublicKey)
		addr, _, _, err := (*(ta.App)).RegisterKey(pubKeyBytes, "secp256k1")
		if err != nil {
			panic(err)
		}

		twitterProfile := &db.TwitterProfile{
			ID:        twitterUser.ID,
			Address:   addr.String(),
			Username:  twitterUser.ScreenName,
			FullName:  twitterUser.Name,
			AvatarURI: twitterUser.ProfileImageURL,
		}
		err = ta.DBClient.UpsertTwitterProfile(twitterProfile)
		if err != nil {
			panic(err)
		}

		// Saves and excrypts the context in the cookie
		hashKey, err := hex.DecodeString(os.Getenv("COOKIE_HASH_KEY"))
		if err != nil {
			panic(err)
		}
		blockKey, err := hex.DecodeString(os.Getenv("COOKIE_ENCRYPT_KEY"))
		if err != nil {
			panic(err)
		}
		s := securecookie.New(hashKey, blockKey)
		cookieValue := map[string]string{
			"twitter-profile-id": twitterUser.IDStr,
			"address":            twitterProfile.Address,
		}
		encodedValue, err := s.Encode("tru-user", cookieValue)
		if err != nil {
			panic(err)
		}

		cookie := http.Cookie{
			Name:     "tru-user",
			HttpOnly: true,
			Value:    encodedValue,
			Expires:  time.Now().Add(2 * time.Hour),
			Domain:   os.Getenv("COOKIE_HOST"),
		}
		http.SetCookie(w, &cookie)
		fmt.Printf("cookie: %v\n\n", os.Getenv("COOKIE_HOST"))
		fmt.Printf("api key: %v\n\n", os.Getenv("TWITTER_API_KEY"))
		http.Redirect(w, req, os.Getenv("OAUTH_SUCCESS_REDIR"), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
