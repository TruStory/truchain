package argument

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath         = "arguments"
	QueryArgumentByID = "id"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryArgumentByID:
			return k.queryArgumentByID(ctx, req)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

func (k Keeper) queryArgumentByID(ctx sdk.Context, req abci.RequestQuery) (res []byte, err sdk.Error) {
	params := app.QueryByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	argument, err := k.Argument(ctx, params.ID)
	if err != nil {
		return
	}

	return app.MustMarshal(argument), nil
}
