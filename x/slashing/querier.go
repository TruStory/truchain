package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QuerySlash                  = "slash"
	QuerySlashes                = "slashes"
	QueryArgumentSlashes        = "argument_slashes"
	QueryArgumentSlasherSlashes = "argument_slasher_slashes"
)

// QuerySlashParams are params for querying slashes by id queries
type QuerySlashParams struct {
	ID uint64 `json:"id"`
}

// QueryArgumentSlashesParams are params for querying slashes by argument id
type QueryArgumentSlashesParams struct {
	ArgumentID uint64 `json:"argument_id"`
}

// QueryArgumentSlashesParams are params for querying slashes by argument id and slasher
type QueryArgumentSlasherSlashesParams struct {
	ArgumentID uint64         `json:"argument_id"`
	Slasher    sdk.AccAddress `json:"slasher"`
}

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QuerySlash:
			return querySlash(ctx, request, keeper)
		case QuerySlashes:
			return querySlashes(ctx, keeper)
		case QueryArgumentSlashes:
			return queryArgumentSlashes(ctx, request, keeper)
		case QueryArgumentSlasherSlashes:
			return queryArgumentSlasherSlashes(ctx, request, keeper)
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

func queryArgumentSlashes(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryArgumentSlashesParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	slashes := k.ArgumentSlashes(ctx, params.ArgumentID)
	bz, jsonErr := k.codec.MarshalJSON(slashes)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return bz, nil
}

func queryArgumentSlasherSlashes(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryArgumentSlasherSlashesParams{}
	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	slashes := k.ArgumentSlasherSlashes(ctx, params.Slasher, params.ArgumentID)
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
