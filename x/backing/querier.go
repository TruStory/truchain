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
	QueryByStoryIDAndCreator    = "storyIDAndCreator"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBackingByID:
			return queryBackingByID(ctx, req, k)
		case QueryBackingAmountByStoryID:
			return queryBackingAmountByStoryID(ctx, req, k)
		case QueryByStoryIDAndCreator:
			return queryByStoryIDAndCreator(ctx, req, k)
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

func queryBackingAmountByStoryID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
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

func queryByStoryIDAndCreator(
	ctx sdk.Context,
	req abci.RequestQuery,
	k ReadKeeper) (res []byte, sdkErr sdk.Error) {

	params := app.QueryByStoryIDAndCreatorParams{}

	sdkErr = app.UnmarshalQueryParams(req, &params)
	if sdkErr != nil {
		return
	}

	// convert address bech32 string to bytes
	addr, err := sdk.AccAddressFromBech32(params.Creator)
	if err != nil {
		return res, sdk.ErrInvalidAddress("Cannot decode address")
	}

	backing, sdkErr := k.BackingByStoryIDAndCreator(ctx, params.StoryID, addr)
	if sdkErr != nil {
		return
	}

	return app.MustMarshal(backing), nil
}
