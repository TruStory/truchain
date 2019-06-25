package claim

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints
const (
	QueryClaim            = "claim"
	QueryClaims           = "claims"
	QueryCommunityClaims  = "community_claims"
	QueryCreatorClaims    = "creator_claims"
	QueryClaimsIDRange    = "claims_id_range"
	QueryClaimsBeforeTime = "claims_before_time"
	QueryClaimsAfterTime  = "claims_after_time"
)

// QueryClaimParams for a single claim
type QueryClaimParams struct {
	ID uint64 `json:"id"`
}

// QueryCommunityClaimsParams for community claims
type QueryCommunityClaimsParams struct {
	CommunityID uint64 `json:"community_id"`
}

// QueryCreatorClaimsParams for community claims
type QueryCreatorClaimsParams struct {
	Creator sdk.AccAddress `json:"creator"`
}

// QueryClaimsIDRangeParams for claims by an id range
type QueryClaimsIDRangeParams struct {
	StartID uint64 `json:"start_id"`
	EndID   uint64 `json:"end_id"`
}

// QueryClaimsTimeParams for claims by time
type QueryClaimsTimeParams struct {
	CreatedTime time.Time `json:"created_time"`
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryClaim:
			return queryClaim(ctx, req, keeper)
		case QueryClaims:
			return queryClaims(ctx, req, keeper)
		case QueryCommunityClaims:
			return queryCommunityClaims(ctx, req, keeper)
		case QueryCreatorClaims:
			return queryCreatorClaims(ctx, req, keeper)
		case QueryClaimsIDRange:
			return queryClaimsByIDRange(ctx, req, keeper)
		case QueryClaimsBeforeTime:
			return queryClaimsBeforeTime(ctx, req, keeper)
		case QueryClaimsAfterTime:
			return queryClaimsAfterTime(ctx, req, keeper)
		}

		return nil, sdk.ErrUnknownRequest("Unknown claim query endpoint")
	}
}

func queryClaim(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryClaimParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}

	claim, ok := keeper.Claim(ctx, params.ID)
	if !ok {
		return nil, ErrUnknownClaim(claim.ID)
	}

	return mustMarshal(claim)
}

func queryClaims(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	claims := keeper.Claims(ctx)

	return mustMarshal(claims)
}

func queryCommunityClaims(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryCommunityClaimsParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	claims := keeper.CommunityClaims(ctx, params.CommunityID)

	return mustMarshal(claims)
}

func queryCreatorClaims(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryCreatorClaimsParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	claims := keeper.CreatorClaims(ctx, params.Creator)

	return mustMarshal(claims)
}

func queryClaimsByIDRange(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryClaimsIDRangeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	claims := keeper.ClaimsBetweenIDs(ctx, params.StartID, params.EndID)

	return mustMarshal(claims)
}

func queryClaimsBeforeTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryClaimsTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	claims := keeper.ClaimsBeforeTime(ctx, params.CreatedTime)

	return mustMarshal(claims)
}

func queryClaimsAfterTime(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryClaimsTimeParams
	codecErr := ModuleCodec.UnmarshalJSON(req.Data, &params)
	if codecErr != nil {
		return nil, ErrJSONParse(codecErr)
	}
	claims := keeper.ClaimsAfterTime(ctx, params.CreatedTime)

	return mustMarshal(claims)
}

func mustMarshal(v interface{}) (result []byte, err sdk.Error) {
	result, jsonErr := codec.MarshalJSONIndent(ModuleCodec, v)
	if jsonErr != nil {
		return nil, ErrJSONParse(jsonErr)
	}

	return
}
