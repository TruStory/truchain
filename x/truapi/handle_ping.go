package truapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
)

// PingRequest is an empty JSON request
type PingRequest struct {}

// PingResponse is a JSON response body representing the result of Ping
type PingResponse struct {
	Pong	bool    `json:"pong"`
}

// HandlePing takes a `PingRequest` and returns a `PingResponse`
func (ta *TruAPI) HandlePing(r *http.Request) chttp.Response {
	rr := new(PingRequest)
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	err = json.Unmarshal(reqBytes, &rr)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	responseBytes, _ := json.Marshal(PingResponse{
		Pong: true,
	})

	return chttp.SimpleResponse(200, responseBytes)
}
