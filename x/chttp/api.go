package chttp

import (
	"context"
	"fmt"
	"net/http"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	abci "github.com/tendermint/tendermint/abci/types"
	tcmn "github.com/tendermint/tendermint/libs/common"
	trpctypes "github.com/tendermint/tendermint/rpc/core/types"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

// MsgTypes is a map of `Msg` type names to empty instances
type MsgTypes map[string]interface{}

// App is implemented by a Cosmos app to provide chain functionality to the API
type App interface {
	RegisterKey(tcmn.HexBytes, string) (sdk.AccAddress, uint64, sdk.Coins, error)
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

// Subrouter returns a mux subrouter.
func (a *API) Subrouter(path string) *mux.Router {
	return a.router.PathPrefix(path).Subrouter()
}

// PathPrefix adds a http.Handler to a path prefix
func (a *API) PathPrefix(path string, handler http.Handler) {
	a.router.PathPrefix(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

// Handle registes a http.Handler
func (a *API) Handle(path string, handler http.Handler) {
	a.router.Handle(path, handler)
}

// Use registers a middleware on the API router
func (a *API) Use(mw func(http.Handler) http.Handler) {
	a.router.Use(mw)
}

// ListenAndServe serves HTTP using the API router
func (a *API) ListenAndServe(addr string) error {
	letsEncryptEnabled := os.Getenv("CHAIN_LETS_ENCRYPT_ENABLED") == "true"
	if !letsEncryptEnabled {
		return http.ListenAndServe(addr, a.router)
	}
	return a.listenAndServe()
}

func (a *API) listenAndServe() error {
	host := os.Getenv("CHAIN_HOST")
	certDir := os.Getenv("CHAIN_LETS_ENCRYPT_CACHE_DIR")
	if certDir == "" {
		certDir = "certs"
	}
	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
	}
	httpServer := &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(redirectHandler),
	}
	secureServer := &http.Server{
		Addr:      ":https",
		Handler:   a.router,
		TLSConfig: m.TLSConfig(),
	}

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		err := httpServer.ListenAndServe()
		if err != nil {
			ctx.Err()
		}
		return err

	})
	g.Go(func() error {
		err := secureServer.ListenAndServeTLS("", "")
		if err != nil {
			ctx.Err()
		}
		return err
	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			httpServer.Shutdown(ctx)
			secureServer.Shutdown(ctx)
			return nil
		}
	})
	return g.Wait()

}

// RunQuery dispatches a query (path + params) to the Cosmos app
func (a *API) RunQuery(path string, params interface{}) abci.ResponseQuery {
	return (*(a.App)).RunQuery(path, params)
}

// DeliverPresigned dispatches a pre-signed query to the Cosmos app
func (a *API) DeliverPresigned(tx auth.StdTx) (*trpctypes.ResultBroadcastTxCommit, error) {
	return (*(a.App)).DeliverPresigned(tx)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
