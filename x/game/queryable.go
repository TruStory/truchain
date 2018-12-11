package game

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath                       = "games"
	QueryGameByID                   = "id"
	QueryChallengeThresholdByGameID = "challengeThresholdByGameID"
)

// QueryGameByIDParams are params for stories by category queries
type QueryGameByIDParams struct {
	ID int64
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(
	gameKeeper ReadKeeper, backingKeeper backing.ReadKeeper) sdk.Querier {

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGameByID:
			return queryGameByID(ctx, req, gameKeeper)
		case QueryChallengeThresholdByGameID:
			return queryChallengeThresholdByGameID(ctx, req, gameKeeper, backingKeeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

// ============================================================================

func queryGameByID(
	ctx sdk.Context, req abci.RequestQuery, gameKeeper ReadKeeper) (
	res []byte, err sdk.Error) {

	params := QueryGameByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	game, err := gameKeeper.Game(ctx, params.ID)
	if err != nil {
		return
	}

	return app.MustMarshal(game), nil
}

func queryChallengeThresholdByGameID(
	ctx sdk.Context,
	req abci.RequestQuery,
	gameKeeper ReadKeeper,
	backingKeeper backing.ReadKeeper) (res []byte, err sdk.Error) {

	params := app.QueryByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	game, err := gameKeeper.Game(ctx, params.ID)
	if err != nil {
		return
	}

	// get the total of all backings on story
	totalBackingAmount, err := backingKeeper.TotalBackingAmount(ctx, game.StoryID)
	if err != nil {
		return nil, err
	}

	challengeThresholdAmount := gameKeeper.ChallengeThreshold(totalBackingAmount)

	return app.MustMarshal(challengeThresholdAmount), nil
}
