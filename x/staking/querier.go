package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	app "github.com/TruStory/truchain/types"
)

const (
	QueryClaimArguments   = "claim_arguments"
	QueryUserArguments    = "user_arguments"
	QueryArgumentStakes   = "argument_stakes"
	QueryUserStakes       = "user_stakes"
	QueryClaimTopArgument = "claim_top_argument"
	QueryEarnedCoins      = "earned_coins"
	QueryTotalEarnedCoins = "total_earned_coins"
)

type QueryClaimArgumentsParams struct {
	ClaimID uint64 `json:"claim_id"`
}

type QueryUserArgumentsParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryArgumentStakesParams struct {
	ArgumentID uint64 `json:"argument_id"`
}

type QueryUserStakesParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryClaimTopArgumentParams struct {
	ClaimID uint64 `json:"claim_id"`
}

type QueryEarnedCoinsParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryTotalEarnedCoinsParams struct {
	Address sdk.AccAddress `json:"address"`
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryClaimArguments:
			return queryClaimArguments(ctx, req, keeper)
		case QueryUserArguments:
			return queryUserArguments(ctx, req, keeper)
		case QueryArgumentStakes:
			return queryArgumentStakes(ctx, req, keeper)
		case QueryUserStakes:
			return queryUserStakes(ctx, req, keeper)
		case QueryClaimTopArgument:
			return queryClaimTopArgument(ctx, req, keeper)
		case QueryEarnedCoins:
			return queryEarnedCoins(ctx, req, keeper)
		case QueryTotalEarnedCoins:
			return queryTotalEarnedCoins(ctx, req, keeper)
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
	arguments := keeper.UserArguments(ctx, params.Address)
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

func queryArgumentStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryArgumentStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.ArgumentStakes(ctx, params.ArgumentID)
	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryUserStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryUserStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.UserStakes(ctx, params.Address)
	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryClaimTopArgument(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryClaimTopArgumentParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	arguments := keeper.ClaimArguments(ctx, params.ClaimID)
	topArgument := Argument{}
	if len(arguments) == 0 {
		bz, err := keeper.codec.MarshalJSON(topArgument)
		if err != nil {
			return nil, ErrJSONParse(err)
		}
		return bz, nil
	}
	for _, a := range arguments {
		if topArgument.ID == 0 {
			topArgument = a
		}
		if topArgument.UpvotedStake.IsLT(a.UpvotedStake) {
			topArgument = a
		}
	}
	bz, err := keeper.codec.MarshalJSON(topArgument)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryEarnedCoins(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryEarnedCoinsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	earnedCoins := keeper.getEarnedCoins(ctx, params.Address)
	bz, err := keeper.codec.MarshalJSON(earnedCoins)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryTotalEarnedCoins(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTotalEarnedCoinsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	earnedCoins := keeper.getEarnedCoins(ctx, params.Address)
	total := sdk.NewInt(0)
	for _, e := range earnedCoins {
		total = total.Add(e.Amount)
	}
	totalStakeEarned := sdk.NewCoin(app.StakeDenom, total)
	bz, err := keeper.codec.MarshalJSON(totalStakeEarned)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}
