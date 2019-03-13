package truapi

import (
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/TruStory/truchain/x/db"
	"github.com/dghubble/gologin/twitter"
	"github.com/gorilla/securecookie"
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

		// NOTE: DATABASE TRANSACTION COULD BE USED IN HERE

		// Fetch the user, if already exists
		currentTwitterProfile, err := ta.DBClient.TwitterProfileByID(twitterUser.ID)
		if err != nil {
			panic(err)
		}
		// if user exists,
		var addr string
		if currentTwitterProfile.ID != 0 {
			addr = currentTwitterProfile.Address
		}

		// Fetch keypair of the user, if already exists
		keyPair, err := ta.DBClient.KeyPairByTwitterProfileID(twitterUser.ID)
		if err != nil {
			panic(err)
		}

		// If not available, create new
		if keyPair.ID == 0 {
			// privKey, err := ecdsa.GenerateKey(ethsecp.S256(), rand.Reader)
			// if err != nil {
			// 	panic(err)
			// }

			keyPair := &db.KeyPair{
				TwitterProfileID: twitterUser.ID,
				// PrivateKey:       fmt.Sprintf("%x", privKey.D),
				// PublicKey:        fmt.Sprintf("%x", privKey.PubKey().Bytes()),
				PrivateKey: "274CB4377D10823E7BE92307FC9A3FD6F28CCBAB9195405EBC3382F5D378A6F7",
				PublicKey:  "02CD2B5105FBAE523891AEF5110DE4D50B79480F81A0CEAEA3A7C844B8451AE213",
			}
			err = ta.DBClient.Add(keyPair)
			if err != nil {
				panic(err)
			}

			// Register with cosmos only if it wasn't registered before.
			if currentTwitterProfile.ID == 0 {
				pubKeyBytes, _ := hex.DecodeString(keyPair.PublicKey)
				newAddr, _, _, err := (*(ta.App)).RegisterKey(pubKeyBytes, "secp256k1")
				if err != nil {
					panic(err)
				}
				addr = newAddr.String()
			}
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
		http.Redirect(w, req, os.Getenv("OAUTH_SUCCESS_REDIR"), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
