package backing

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath                   = "backings"
	QueryBackingByID            = "id"
	QueryBackingAmountByStoryID = "totalAmountByStoryID"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBackingByID:
			return queryBackingByID(ctx, req, k)
		case QueryBackingAmountByStoryID:
			return queryBackinAmountByStoryID(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

// ============================================================================

func queryBackingByID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	params := app.QueryByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	backing, err := k.Backing(ctx, params.ID)
	if err != nil {
		return
	}

	return app.MustMarshal(backing), nil
}

func queryBackinAmountByStoryID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	params := app.QueryByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	backingTotal, err := k.TotalBackingAmount(ctx, params.ID)
	if err != nil {
		return
	}

	return app.MustMarshal(backingTotal), nil
}
