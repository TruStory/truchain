package claim

import (
	"encoding/json"
	"strings"
	"testing"

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
	cdcErr := moduleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 0, len(claims))
}

func TestQueryClaims(t *testing.T) {
	ctx, keeper := mockDB()

	fakeClaim(ctx, keeper)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryClaims}, "/"),
		Data: []byte{},
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryClaims}, query)
	require.NoError(t, err)

	var claims []Claim
	cdcErr := moduleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 1, len(claims))
}

func TestQueryCommunityClaims(t *testing.T) {
	ctx, keeper := mockDB()

	fakeClaim(ctx, keeper)

	queryParams := QueryCommunityClaimsParams{
		CommunityID: uint64(1),
	}
	queryParamsBytes, jsonErr := json.Marshal(queryParams)
	require.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryCommunityClaims}, "/"),
		Data: queryParamsBytes,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryCommunityClaims}, query)
	require.NoError(t, err)

	var claims []Claim
	cdcErr := moduleCodec.UnmarshalJSON(resBytes, &claims)
	require.NoError(t, cdcErr)
	require.Equal(t, 1, len(claims))
}
