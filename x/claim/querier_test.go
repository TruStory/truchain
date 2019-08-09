package claim

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

const custom = "custom"

func TestQueryClaims_NoneFound(t *testing.T) {
	ctx, keeper := mockDB()

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryClaims}, "/"),
		Data: []byte{},
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryClaims}, query)
	require.NoError(t, err)

	var claims []Claim
	cdcErr := ModuleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 0, len(claims))
}

func TestQueryClaims(t *testing.T) {
	ctx, keeper := mockDB()

	fakeClaim(ctx, keeper, "crypto")

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryClaims}, "/"),
		Data: []byte{},
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryClaims}, query)
	require.NoError(t, err)

	var claims []Claim
	cdcErr := ModuleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 1, len(claims))
}

func TestQueryClaimsByIDs(t *testing.T) {
	ctx, keeper := mockDB()

	fakeClaim(ctx, keeper, "crypto")
	fakeClaim(ctx, keeper, "crypto")
	fakeClaim(ctx, keeper, "crypto")
	fakeClaim(ctx, keeper, "crypto")
	fakeClaim(ctx, keeper, "crypto")

	queryParams := QueryClaimsParams{
		IDs: []uint64{1, 3, 4},
	}
	queryParamsBytes, jsonErr := ModuleCodec.MarshalJSON(queryParams)
	require.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryClaimsByIDs}, "/"),
		Data: queryParamsBytes,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryClaimsByIDs}, query)
	require.NoError(t, err)

	var claims Claims
	cdcErr := ModuleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 3, len(claims))
}

func TestQueryCommunityClaims(t *testing.T) {
	ctx, keeper := mockDB()

	fakeClaim(ctx, keeper, "crypto")

	queryParams := QueryCommunityClaimsParams{
		CommunityID: "crypto",
	}
	queryParamsBytes, jsonErr := ModuleCodec.MarshalJSON(queryParams)
	require.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryCommunityClaims}, "/"),
		Data: queryParamsBytes,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryCommunityClaims}, query)
	require.NoError(t, err)

	var claims []Claim
	cdcErr := ModuleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 1, len(claims))
}

func TestQueryCommunitiesClaims(t *testing.T) {
	ctx, keeper := mockDB()

	fakeClaim(ctx, keeper, "crypto")
	fakeClaim(ctx, keeper, "meme")

	queryParams := QueryCommunitiesClaimsParams{
		CommunityIDs: []string{"crypto", "meme"},
	}
	queryParamsBytes, jsonErr := ModuleCodec.MarshalJSON(queryParams)
	require.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryCommunitiesClaims}, "/"),
		Data: queryParamsBytes,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryCommunitiesClaims}, query)
	require.NoError(t, err)

	var claims []Claim
	cdcErr := ModuleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 2, len(claims))
}

func TestQueryCreatorClaims(t *testing.T) {
	ctx, keeper := mockDB()

	claim := fakeClaim(ctx, keeper, "crypto")

	queryParams := QueryCreatorClaimsParams{
		Creator: claim.Creator,
	}
	queryParamsBytes, jsonErr := ModuleCodec.MarshalJSON(queryParams)
	require.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryCreatorClaims}, "/"),
		Data: queryParamsBytes,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryCreatorClaims}, query)
	require.NoError(t, err)

	var claims []Claim
	cdcErr := ModuleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 1, len(claims))
}

func TestQueryParams_Success(t *testing.T) {
	ctx, keeper := mockDB()

	onChainParams := keeper.GetParams(ctx)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryParams),
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryParams}, query)
	assert.Nil(t, err)

	var returnedParams Params
	sdkErr := ModuleCodec.UnmarshalJSON(resBytes, &returnedParams)
	assert.Nil(t, sdkErr)
	assert.Equal(t, returnedParams, onChainParams)
}
