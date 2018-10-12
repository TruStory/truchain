package story

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryCategoryStories = "category/stories"
)

// QueryCategoryStoriesParams are params for query 'category/stories'
type QueryCategoryStoriesParams struct {
	CategoryID int64
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryCategoryStories:
			return queryStoriesWithCategory(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

// ============================================================================

func queryStoriesWithCategory(
	ctx sdk.Context,
	req abci.RequestQuery,
	k ReadKeeper) (res []byte, sdkErr sdk.Error) {

	// deserialize query params
	var params QueryCategoryStoriesParams
	err := k.GetCodec().UnmarshalJSON(req.Data, &params)
	if err != nil {
		return res,
			sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
	}

	// fetch stories
	stories, sdkErr := k.GetStoriesWithCategory(ctx, params.CategoryID)
	if sdkErr != nil {
		return res, sdkErr
	}

	// serialize into pretty JSON bytes
	res, err = codec.MarshalJSONIndent(k.GetCodec(), stories)
	if err != nil {
		panic("Could not marshal result to JSON")
	}

	return
}
