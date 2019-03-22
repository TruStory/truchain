package truapi

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"
	"github.com/btcsuite/btcd/btcec"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// RegistrationRequest is a JSON request body representing a twitter profile that a user wishes to register
type RegistrationRequest struct {
	AuthToken       string `json:"auth_token"`
	AuthTokenSecret string `json:"auth_token_secret"`
}

// RegistrationResponse is a JSON response body representing the result of registering a key
type RegistrationResponse struct {
	UserID               string `json:"userId"`
	Username             string `json:"username"`
	Fullname             string `json:"fullname"`
	Address              string `json:"address"`
	AuthenticationCookie string `json:"authenticationCookie"`
}

// HandleRegistration takes a `RegistrationRequest` and returns a `RegistrationResponse`
func (ta *TruAPI) HandleRegistration(r *http.Request) chttp.Response {
	rr := new(RegistrationRequest)
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	err = json.Unmarshal(reqBytes, &rr)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	// Get the Twitter User from the auth token
	twitterUser, err := getTwitterUser(rr.AuthToken, rr.AuthTokenSecret)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	return RegisterTwitterUser(ta, twitterUser)
}

// RegisterTwitterUser registers a new twitter user
func RegisterTwitterUser(ta *TruAPI, twitterUser *twitter.User) chttp.Response {
	addr, err := CalibrateUser(ta, twitterUser)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	twitterProfile := &db.TwitterProfile{
		ID:        twitterUser.ID,
		Address:   addr,
		Username:  twitterUser.ScreenName,
		FullName:  twitterUser.Name,
		AvatarURI: strings.Replace(twitterUser.ProfileImageURL, "_normal", "_bigger", 1),
	}

	err = ta.DBClient.UpsertTwitterProfile(twitterProfile)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	cookieValue, err := cookies.MakeLoginCookieValue(twitterProfile)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	responseBytes, _ := json.Marshal(RegistrationResponse{
		UserID:               twitterUser.IDStr,
		Username:             twitterUser.ScreenName,
		Fullname:             twitterUser.Name,
		Address:              addr,
		AuthenticationCookie: cookieValue,
	})

	return chttp.SimpleResponse(201, responseBytes)
}

// CalibrateUser takes a twitter authenticated user and makes sure it has properly
// been calibrated in the database with all proper keypairs
func CalibrateUser(ta *TruAPI, twitterUser *twitter.User) (string, error) {
	currentTwitterProfile, err := ta.DBClient.TwitterProfileByID(twitterUser.ID)
	if err != nil {
		return "", err
	}
	// if user exists,
	var addr string
	if currentTwitterProfile.ID != 0 {
		addr = currentTwitterProfile.Address
	}

	// Fetch keypair of the user, if already exists
	keyPair, err := ta.DBClient.KeyPairByTwitterProfileID(twitterUser.ID)
	if err != nil {
		return "", err
	}

	// If not available, create new
	if keyPair.ID == 0 {
		newKeyPair, _ := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			return "", err
		}
		// We are converting the private key of the new key pair in hex string,
		// then back to byte slice, and finally regenerating the private (suppressed) and public key from it.
		// This way, it returns the kind of public key that cosmos understands.
		_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), []byte(fmt.Sprintf("%x", newKeyPair.Serialize())))

		keyPair := &db.KeyPair{
			TwitterProfileID: twitterUser.ID,
			PrivateKey:       fmt.Sprintf("%x", newKeyPair.Serialize()),
			PublicKey:        fmt.Sprintf("%x", pubKey.SerializeCompressed()),
		}
		err = ta.DBClient.Add(keyPair)
		if err != nil {
			return "", err
		}

		// Register with cosmos only if it wasn't registered before.
		if currentTwitterProfile.ID == 0 {
			pubKeyBytes, _ := hex.DecodeString(keyPair.PublicKey)
			newAddr, _, _, err := (*(ta.App)).RegisterKey(pubKeyBytes, "secp256k1")
			if err != nil {
				return "", err
			}
			addr = newAddr.String()
		}
	}

	return addr, nil
}

func getTwitterUser(authToken string, authTokenSecret string) (*twitter.User, error) {
	ctx := context.Background()
	config := oauth1.NewConfig(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_SECRET"))

	httpClient := config.Client(ctx, oauth1.NewToken(authToken, authTokenSecret))
	twitterClient := twitter.NewClient(httpClient)
	accountVerifyParams := &twitter.AccountVerifyParams{
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
		IncludeEmail:    twitter.Bool(false),
	}
	user, _, err := twitterClient.Accounts.VerifyCredentials(accountVerifyParams)
	if err != nil {
		return nil, err
	}

	return user, nil
}
