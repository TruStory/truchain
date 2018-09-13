package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/lcd"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	return lcd.ServeCommand(cdc)
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	// GET /stories
	r.HandleFunc("/stories",
		GetStoriesHandlerFn(cdc, kb, ctx)).Methods("GET")
}

// GetStoriesHandlerFn - http request handler to get stories
func GetStoriesHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// body, err := ioutil.ReadAll(r.Body)
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		key := []byte("stories")
		res, err := ctx.QueryStore(key, "stories")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// the query will return empty if there is no data
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		output, err := json.MarshalIndent(res, "", " ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)
	}
}
