package truapi

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/gorilla/securecookie"
)

// HandleUnsigned takes a `HandleUnsignedRequest` and returns a `HandleUnsignedResponse`
func (ta *TruAPI) HandleUnsigned(r *http.Request) chttp.Response {
	txr := new(chttp.UnsignedRequest)
	jsonBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}

	err = json.Unmarshal(jsonBytes, txr)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	// Get the user context
	truUser, err := GetTruUserFromCookie(r)
	if err == http.ErrNoCookie {
		return chttp.SimpleErrorResponse(401, err)
	}
	if err != nil {
		panic(err)
	}

	twitterProfileID, err := strconv.ParseInt(truUser["twitter-profile-id"], 10, 64)
	if err != nil {
		panic(err)
	}

	// Fetch keypair of the user
	keyPair, err := ta.DBClient.KeyPairByTwitterProfileID(twitterProfileID)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}
	if keyPair.ID == 0 {
		// keypair doesn't exist
		return chttp.SimpleErrorResponse(400, errors.New("keypair does not exist on the server"))
	}

	tx, err := ta.NewUnsignedStdTx(*txr, keyPair)

	if err != nil {
		fmt.Println("Error decoding tx: ", err)
		return chttp.SimpleErrorResponse(400, err)
	}

	res, err := ta.DeliverPresigned(tx)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	resBytes, _ := json.Marshal(res)

	return chttp.SimpleResponse(200, resBytes)
}

// GetTruUserFromCookie gets the user context from the request's cookie
func GetTruUserFromCookie(r *http.Request) (map[string]string, error) {
	truUser, err := r.Cookie("tru-user")
	if err != nil {
		return nil, err
	}

	hashKey, err := hex.DecodeString(os.Getenv("COOKIE_HASH_KEY"))
	if err != nil {
		return nil, err
	}
	blockKey, err := hex.DecodeString(os.Getenv("COOKIE_ENCRYPT_KEY"))
	if err != nil {
		return nil, err
	}
	var s = securecookie.New(hashKey, blockKey)

	decodedTruUser := make(map[string]string)
	err = s.Decode("tru-user", truUser.Value, &decodedTruUser)
	if err != nil {
		return nil, err
	}

	return decodedTruUser, nil
}
