package params

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath = "params"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		return queryParams(ctx, req, k)
	}
}

// ============================================================================

func queryParams(ctx sdk.Context, req abci.RequestQuery, k Keeper) (res []byte, err sdk.Error) {
	return mustMarshal(k.Params(ctx))
}

func mustMarshal(v interface{}) (result []byte, err sdk.Error) {
	result, jsonErr := codec.MarshalJSONIndent(ModuleCodec, v)
	if jsonErr != nil {
		return nil, ErrJSONParse(jsonErr)
	}

	return
}
