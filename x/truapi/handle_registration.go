package truapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tcmn "github.com/tendermint/tendermint/libs/common"
)

type RegistrationRequest struct {
	PubKeyAlgo string        `json:"pubkey_algo"`
	PubKey     tcmn.HexBytes `json:"pubkey"`
}

type RegistrationResponse struct {
	Address       string    `json:"address"`
	AccountNumber int64     `json:"account_number"`
	Sequence      int64     `json:"sequence"`
	Coins         sdk.Coins `json:"coins"`
}

func (ta *TruApi) HandleRegistration(r *http.Request) chttp.Response {
	rr := new(RegistrationRequest)
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	err = json.Unmarshal(reqBytes, &rr)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	addr, num, coins, err := (*(ta.App)).RegisterKey(rr.PubKey, rr.PubKeyAlgo)

	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	responseBytes, _ := json.Marshal(RegistrationResponse{
		Address:       addr.String(),
		AccountNumber: num,
		Sequence:      0,
		Coins:         coins,
	})

	return chttp.SimpleResponse(201, responseBytes)
}
