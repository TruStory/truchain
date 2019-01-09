package vote

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath                = "votes"
	QueryByGameID            = "gameID"
	QueryByStoryIDAndCreator = "storyIDAndCreator"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryByGameID:
			return queryByGameID(ctx, req, k)
		case QueryByStoryIDAndCreator:
			return queryByStoryIDAndCreator(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown query endpoint")
		}
	}
}

// ============================================================================

func queryByGameID(
	ctx sdk.Context,
	req abci.RequestQuery,
	k ReadKeeper) (res []byte, sdkErr sdk.Error) {

	params := app.QueryByIDParams{}

	sdkErr = app.UnmarshalQueryParams(req, &params)
	if sdkErr != nil {
		return
	}

	tokenVotes, sdkErr := k.TokenVotesByGameID(ctx, params.ID)
	if sdkErr != nil {
		return
	}

	return app.MustMarshal(tokenVotes), nil
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

	tokenVote, sdkErr := k.TokenVotesByStoryIDAndCreator(ctx, params.StoryID, addr)
	if sdkErr != nil {
		return
	}

	return app.MustMarshal(tokenVote), nil
}
