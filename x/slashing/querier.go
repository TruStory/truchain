package slashing

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QuerySlash   = "id"
	QuerySlashes = "all"
)

// QuerySlashParams are params for querying slashes by id queries
type QuerySlashParams struct {
	ID uint64
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QuerySlash:
			return querySlash(ctx, request, keeper)
		case QuerySlashes:
			return querySlashes(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown truchain query endpoint: slashing/%s", path[0]))
		}
	}
}

func querySlash(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QuerySlashParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	slash, err := k.Slash(ctx, params.ID)
	if err != nil {
		return
	}

	return mustMarshal(slash), nil
}

func querySlashes(ctx sdk.Context, k Keeper) (result []byte, err sdk.Error) {
	slashes := k.Slashes(ctx)
	return mustMarshal(slashes), nil
}

func unmarshalQueryParams(request abci.RequestQuery, params interface{}) (sdkErr sdk.Error) {
	err := json.Unmarshal(request.Data, params)
	if err != nil {
		sdkErr = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
		return
	}
	return
}

func mustMarshal(v interface{}) (result []byte) {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic("Could not marshal result to JSON")
	}
	return
}
