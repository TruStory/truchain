package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QuerySlash        = "slash"
	QuerySlashes      = "slashes"
	QueryStakeSlashes = "stake_slashes"
)

// QuerySlashParams are params for querying slashes by id queries
type QuerySlashParams struct {
	ID uint64 `json:"id"`
}

// QueryStakeSlashesParams are params for querying slashes by stake id
type QueryStakeSlashesParams struct {
	StakeID uint64 `json:"stake_id"`
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QuerySlash:
			return querySlash(ctx, request, keeper)
		case QuerySlashes:
			return querySlashes(ctx, keeper)
		case QueryStakeSlashes:
			return queryStakeSlashes(ctx, request, keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown truchain query endpoint: slashing/%s", path[0]))
		}
	}
}

func querySlash(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QuerySlashParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	slash, err := k.Slash(ctx, params.ID)
	if err != nil {
		return
	}
	bz, jsonErr := k.codec.MarshalJSON(slash)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return bz, nil
}

func querySlashes(ctx sdk.Context, k Keeper) (result []byte, err sdk.Error) {
	slashes := k.Slashes(ctx)
	bz, jsonErr := k.codec.MarshalJSON(slashes)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return bz, nil
}

func queryStakeSlashes(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryStakeSlashesParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	slashes := k.StakeSlashes(ctx, params.StakeID)
	bz, jsonErr := k.codec.MarshalJSON(slashes)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return bz, nil
}

func unmarshalQueryParams(request abci.RequestQuery, params interface{}) (sdkErr sdk.Error) {
	err := ModuleCodec.UnmarshalJSON(request.Data, params)
	if err != nil {
		sdkErr = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
		return
	}
	return
}
