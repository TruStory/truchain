package truapi


import (
	"encoding/json"
  "errors"
  "fmt"
	"io/ioutil"
	"net/http"
  "github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/chttp"
)

// DeviceResponse is a JSON response body representing the result of updating the device Token

type DeviceResponse struct {
	Address       string    `json:"address"`
  Token       string    `json:"token"`
}

type DeviceRequest struct {
	Address       string    `json:"address"`
  Token       string    `json:"token"`
}

// HandleDevice takes a `DeviceRequest` and returns a `DeviceResponse`

func (ta *TruAPI) HandleDevice(r *http.Request) chttp.Response {
  fmt.Println("HEYYYYYY")
	rr := new(DeviceRequest)
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	err = json.Unmarshal(reqBytes, &rr)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	if rr.Token == "" {
		return chttp.SimpleErrorResponse(400, errors.New("Device Token is required"))
	}

	deviceToken := &db.DeviceToken{
		Address:  rr.Address,
    Token: rr.Token,
	}

	err = ta.DBClient.InsertDeviceToken(deviceToken)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	responseBytes, _ := json.Marshal(DeviceResponse{
		Address: rr.Address,
    Token: rr.Token,
	})

	return chttp.SimpleResponse(201, responseBytes)
}
