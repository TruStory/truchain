package category

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath           = "categories"
	QueryCategoriesByID = "id"
)

// QueryCategoryParams are params for  by category queries
type QueryCategoryParams struct {
	ID string
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryCategoriesByID:
			return queryCategoryByID(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint: categories/" + path[0])
		}
	}
}

// ============================================================================

func queryCategoryByID(
	ctx sdk.Context,
	req abci.RequestQuery,
	k ReadKeeper) (res []byte, err sdk.Error) {

	// get query params
	params, err := unmarshalQueryParams(k, req)
	if err != nil {
		return
	}

	cid, parseErr := strconv.ParseInt(params.ID, 10, 64)

	if parseErr != nil {
		return res,
			sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
	}

	// fetch
	category, sdkErr := k.GetCategory(ctx, cid)
	if sdkErr != nil {
		return res, sdkErr
	}

	// return JSON bytes
	return marshalCategory(k, category)
}

// unmarshal query params into struct
func unmarshalQueryParams(
	k ReadKeeper,
	req abci.RequestQuery) (params QueryCategoryParams, sdkErr sdk.Error) {
	err := k.GetCodec().UnmarshalJSON(req.Data, &params)
	if err != nil {
		return params,
			sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
	}
	return
}

// marshal  into pretty JSON bytes
func marshalCategory(k ReadKeeper, c Category) (res []byte, sdkErr sdk.Error) {
	res, err := codec.MarshalJSONIndent(k.GetCodec(), c)
	if err != nil {
		panic("Could not marshal result to JSON")
	}
	return
}
