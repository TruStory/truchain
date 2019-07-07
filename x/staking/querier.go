package staking

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryClaimArgument       = "claim_argument"
	QueryClaimArguments      = "claim_arguments"
	QueryUserArguments       = "user_arguments"
	QueryArgumentStakes      = "argument_stakes"
	QueryCommunityStakes     = "community_stakes"
	QueryStake               = "stake"
	QueryArgumentsByIDs      = "arguments_ids"
	QueryUserStakes          = "user_stakes"
	QueryUserCommunityStakes = "user_community_stakes"
	QueryClaimTopArgument    = "claim_top_argument"
	QueryEarnedCoins         = "earned_coins"
	QueryTotalEarnedCoins    = "total_earned_coins"
)

type QueryClaimArgumentParams struct {
	ArgumentID uint64 `json:"argument_id"`
}

type QueryClaimArgumentsParams struct {
	ClaimID uint64 `json:"claim_id"`
}

type QueryUserArgumentsParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryArgumentStakesParams struct {
	ArgumentID uint64 `json:"argument_id"`
}

type QueryCommunityStakesParams struct {
	CommunityID string `json:"community_id"`
}

type QueryStakeParams struct {
	StakeID uint64 `json:"stake_id"`
}

type QueryArgumentsByIDsParams struct {
	ArgumentIDs []uint64 `json:"argument_ids"`
}

type QueryUserStakesParams struct {
	Address sdk.AccAddress `json:"address"`
}

type QueryUserCommunityStakesParams struct {
	Address     sdk.AccAddress `json:"address"`
	CommunityID string         `json:"community_id"`
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
		case QueryClaimArgument:
			return queryClaimArgument(ctx, req, keeper)
		case QueryClaimArguments:
			return queryClaimArguments(ctx, req, keeper)
		case QueryUserArguments:
			return queryUserArguments(ctx, req, keeper)
		case QueryArgumentStakes:
			return queryArgumentStakes(ctx, req, keeper)
		case QueryCommunityStakes:
			return queryCommunityStakes(ctx, req, keeper)
		case QueryStake:
			return queryStake(ctx, req, keeper)
		case QueryArgumentsByIDs:
			return queryArgumentsByIDs(ctx, req, keeper)
		case QueryUserStakes:
			return queryUserStakes(ctx, req, keeper)
		case QueryUserCommunityStakes:
			return queryUserCommunityStakes(ctx, req, keeper)
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

func queryClaimArgument(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryClaimArgumentParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	argument, ok := keeper.getArgument(ctx, params.ArgumentID)
	if !ok {
		return nil, ErrCodeUnknownArgument(params.ArgumentID)
	}
	bz, err := keeper.codec.MarshalJSON(argument)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
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

func queryCommunityStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryCommunityStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.CommunityStakes(ctx, params.CommunityID)
	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryStake(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryStakeParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}

	stakes, ok := keeper.getStake(ctx, params.StakeID)
	if !ok {
		return nil, ErrCodeUnknownStake(params.StakeID)
	}

	bz, err := keeper.codec.MarshalJSON(stakes)
	if err != nil {
		return nil, ErrJSONParse(err)
	}
	return bz, nil
}

func queryArgumentsByIDs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryArgumentsByIDsParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	var arguments []Argument
	for _, id := range params.ArgumentIDs {
		a, ok := keeper.getArgument(ctx, id)
		if !ok {
			return nil, ErrCodeUnknownArgument(id)
		}
		arguments = append(arguments, a)
	}

	bz, err := keeper.codec.MarshalJSON(arguments)
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

func queryUserCommunityStakes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryUserCommunityStakesParams
	err := keeper.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, ErrInvalidQueryParams(err)
	}
	stakes := keeper.UserCommunityStakes(ctx, params.Address, params.CommunityID)
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
