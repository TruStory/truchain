package game

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath     = "games"
	QueryGameByID = "id"
)

// QueryGameByIDParams are params for stories by category queries
type QueryGameByIDParams struct {
	ID int64
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGameByID:
			return queryGameByID(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

// ============================================================================

func queryGameByID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	params := QueryGameByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	game, err := k.Game(ctx, params.ID)
	if err != nil {
		return
	}

	return app.MustMarshal(game), nil
}
