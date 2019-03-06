package voting

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath                 = "voting"
	QueryVoteResultsByStoryID = "queryVoteResultsByStoryID"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryVoteResultsByStoryID:
			return queryVoteResultsByStoryID(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown query endpoint")
		}
	}
}

// ============================================================================

func queryVoteResultsByStoryID(
	ctx sdk.Context,
	req abci.RequestQuery,
	k ReadKeeper) (res []byte, sdkErr sdk.Error) {

	params := app.QueryByIDParams{}

	sdkErr = app.UnmarshalQueryParams(req, &params)
	if sdkErr != nil {
		return
	}

	voteResults, sdkErr := k.GetVoteResultsByStoryID(ctx, params.ID)
	if sdkErr != nil {
		return
	}

	return app.MustMarshal(voteResults), nil
}
