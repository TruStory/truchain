package community

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryCommunity   = "id"
	QueryCommunities = "all"
)

// QueryCommunityParams are params for querying communities by id queries
type QueryCommunityParams struct {
	ID int64
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, request abci.RequestQuery) (result []byte, err sdk.Error) {
		switch path[0] {
		case QueryCommunity:
			return queryCommunity(ctx, request, k)
		case QueryCommunities:
			return queryCommunities(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint: communities/" + path[0])
		}
	}
}

// ============================================================================ //
// QUERIERS BELOW
// ============================================================================ //

func queryCommunity(ctx sdk.Context, request abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	params := QueryCommunityParams{}

	if err = unmarshalQueryParams(request, &params); err != nil {
		return
	}

	community, err := k.Community(ctx, params.ID)

	if err != nil {
		return
	}

	return mustMarshal(community), nil
}

func queryCommunities(ctx sdk.Context, k Keeper) (result []byte, err sdk.Error) {
	communities, err := k.Communities(ctx)

	if err != nil {
		return
	}

	return mustMarshal(communities), nil
}

func unmarshalQueryParams(request abci.RequestQuery, params interface{}) (sdkErr sdk.Error) {
	err := json.Unmarshal(request.Data, params)
	if err != nil {
		sdkErr = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", err.Error()))
		return
	}
	return
}

func mustMarshal(v interface{}) (result []byte) {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic("Could not marshal result to JSON")
	}
	return
}
