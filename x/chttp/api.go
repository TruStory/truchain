package chttp

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	abci "github.com/tendermint/tendermint/abci/types"
	tcmn "github.com/tendermint/tendermint/libs/common"
	trpctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type MsgTypes map[string]interface{} // Map of Msg type names to empty instances

type App interface {
	RegisterKey(tcmn.HexBytes, string) (*sdk.AccAddress, int64, *sdk.Coins, error)
	RunQuery(string, interface{}) abci.ResponseQuery
	DeliverPresigned(auth.StdTx) (*trpctypes.ResultBroadcastTxCommit, error)
}

type Api struct {
	App       *App
	Supported MsgTypes
	router    *mux.Router
}

func NewApi(app *App, supported MsgTypes) *Api {
	a := Api{App: app, Supported: supported, router: mux.NewRouter()}
	return &a
}

func (a *Api) HandleFunc(path string, h Handler) {
	a.router.HandleFunc(path, h.HandlerFunc())
}

func (a *Api) Use(mw func(http.Handler) http.Handler) {
	a.router.Use(mw)
}

func (a *Api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, a.router)
}

func (a *Api) RunQuery(path string, params interface{}) abci.ResponseQuery {
	return (*(a.App)).RunQuery(path, params)
}
