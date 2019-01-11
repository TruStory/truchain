package params

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath = "params"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier() sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		return queryParams(ctx, req)
	}
}

// ============================================================================

func queryParams(ctx sdk.Context, req abci.RequestQuery) (res []byte, err sdk.Error) {
	params := DefaultParams()

	return app.MustMarshal(params), nil
}
