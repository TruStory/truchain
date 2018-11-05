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

// MsgTypes is a map of `Msg` type names to empty instances
type MsgTypes map[string]interface{}

// App is implemented by a Cosmos app to provide chain functionality to the API
type App interface {
	RegisterKey(tcmn.HexBytes, string) (sdk.AccAddress, int64, sdk.Coins, error)
	RunQuery(string, interface{}) abci.ResponseQuery
	DeliverPresigned(auth.StdTx) (*trpctypes.ResultBroadcastTxCommit, error)
}

// API presents the functionality of a Cosmos app over HTTP
type API struct {
	App       *App
	Supported MsgTypes
	router    *mux.Router
}

// NewAPI creates an `API` struct from an `App` and a `MsgTypes` schema
func NewAPI(app *App, supported MsgTypes) *API {
	a := API{App: app, Supported: supported, router: mux.NewRouter()}
	return &a
}

// HandleFunc registers a `chttp.Handler` on the API router
func (a *API) HandleFunc(path string, h Handler) {
	a.router.HandleFunc(path, h.HandlerFunc())
}

// Use registers a middleware on the API router
func (a *API) Use(mw func(http.Handler) http.Handler) {
	a.router.Use(mw)
}

// ListenAndServe serves HTTP using the API router
func (a *API) ListenAndServe(addr string) error {
	go func() { _ = http.ListenAndServe("localhost:3030", nil) }()
	return http.ListenAndServe(addr, a.router)
}

// RunQuery dispatches a query (path + params) to the Cosmos app
func (a *API) RunQuery(path string, params interface{}) abci.ResponseQuery {
	return (*(a.App)).RunQuery(path, params)
}

// DeliverPresigned dispatches a pre-signed query to the Cosmos app
func (a *API) DeliverPresigned(tx auth.StdTx) (*trpctypes.ResultBroadcastTxCommit, error) {
	return (*(a.App)).DeliverPresigned(tx)
}
