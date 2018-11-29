package story

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath                = "stories"
	QueryStoriesByCategoryID = "category"
)

// QueryCategoryStoriesParams are params for stories by category queries
type QueryCategoryStoriesParams struct {
	CategoryID int64
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryStoriesByCategoryID:
			return queryStoriesByCategoryID(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

// ============================================================================

func queryStoriesByCategoryID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	params := QueryCategoryStoriesParams{}

	if err = unmarshalQueryParams(req, &params); err != nil {
		return
	}

	stories, err := k.StoriesByCategoryID(ctx, params.CategoryID)

	if err != nil {
		return
	}

	return mustMarshal(stories), nil
}

func unmarshalQueryParams(req abci.RequestQuery, params interface{}) (sdkErr sdk.Error) {
	parseErr := json.Unmarshal(req.Data, params)
	if parseErr != nil {
		sdkErr = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", parseErr.Error()))
		return
	}
	return
}

func mustMarshal(v interface{}) (res []byte) {
	res, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic("Could not marshal result to JSON")
	}
	return
}
