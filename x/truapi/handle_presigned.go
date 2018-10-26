package truapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
)

// HandlePresigned dispatches a `chttp.PresignedRequest` to a Cosmos app
func (ta *TruAPI) HandlePresigned(r *http.Request) chttp.Response {
	txr := new(chttp.PresignedRequest)
	jsonBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}

	err = json.Unmarshal(jsonBytes, txr)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	tx, err := ta.NewPresignedStdTx(*txr)

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
