package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryClaimArguments = "claim_arguments"
	QueryUserArguments  = "user_arguments"
)

type QueryClaimArgumentsParams struct {
	ClaimID uint64 `json:"claim_id"`
}

type QueryUserArgumentsParams struct {
	address sdk.AccAddress `json:"address"`
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryClaimArguments:
			return queryClaimArguments(ctx, req, keeper)
		case QueryUserArguments:
			return queryUserArguments(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown staking query endpoint")
		}
	}
}

func queryUserArguments(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryUserArgumentsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	arguments := keeper.UserArguments(ctx, params.address)
	bz, err := keeper.codec.MarshalJSON(arguments)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryClaimArguments(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryClaimArgumentsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	arguments := keeper.ClaimArguments(ctx, params.ClaimID)
	bz, err := keeper.codec.MarshalJSON(arguments)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}
