package community

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
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
	ID uint64
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
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown truchain query endpoint: commmunity/%s", path[0]))
		}
	}
}

func queryCommunity(ctx sdk.Context, req abci.RequestQuery, k Keeper) (result []byte, err sdk.Error) {
	var params QueryCommunityParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	community, err := k.Community(ctx, params.ID)
	if err != nil {
		return
	}

	return mustMarshal(community)
}

func queryCommunities(ctx sdk.Context, k Keeper) (result []byte, err sdk.Error) {
	communities := k.Communities(ctx)
	return mustMarshal(communities)
}

func mustMarshal(v interface{}) (result []byte, err sdk.Error) {
	result, jsonErr := codec.MarshalJSONIndent(ModuleCodec, v)
	if jsonErr != nil {
		return nil, ErrJSONParse(jsonErr)
	}

	return
}
