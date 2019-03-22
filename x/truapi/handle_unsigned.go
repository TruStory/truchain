package truapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/cookies"
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

	// Get the authenticated user
	user, err := cookies.GetAuthenticatedUser(r)
	if err == http.ErrNoCookie {
		return chttp.SimpleErrorResponse(401, err)
	}
	if err != nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	// Fetch keypair of the user
	keyPair, err := ta.DBClient.KeyPairByTwitterProfileID(user.TwitterProfileID)
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
