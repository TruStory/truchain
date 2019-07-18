package community

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

const custom = "custom"

func TestQueryCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, id, description := getFakeCommunityParams()
	creator := keeper.GetParams(ctx).CommunityAdmins[0]
	createdCommunity, err := keeper.NewCommunity(ctx, name, id, description, creator)
	assert.Nil(t, err)

	params, jsonErr := ModuleCodec.MarshalJSON(QueryCommunityParams{
		ID: "randomness",
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryCommunity}, "/"),
		Data: params,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryCommunity}, query)
	require.NoError(t, err)

	var returnedCommunity Community
	jsonErr = ModuleCodec.UnmarshalJSON(resBytes, &returnedCommunity)
	assert.Nil(t, jsonErr)
	assert.Equal(t, returnedCommunity.Name, createdCommunity.Name)
	assert.Equal(t, returnedCommunity.ID, createdCommunity.ID)
	assert.Equal(t, returnedCommunity.Description, createdCommunity.Description)
}

func TestQueryCommunity_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	params, err := ModuleCodec.MarshalJSON(QueryCommunityParams{
		ID: "test",
	})
	require.Nil(t, err)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryCommunity}, "/"),
		Data: params,
	}

	querier := NewQuerier(keeper)
	_, sdkErr := querier(ctx, []string{QueryCommunity}, query)

	require.NotNil(t, sdkErr)
	require.Equal(t, ErrCommunityNotFound("test").Code(), sdkErr.Code(), "should get error")
}

func TestQueryCommunities_Success(t *testing.T) {
	ctx, keeper := mockDB()

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryCommunities}, "/"),
		Data: []byte{},
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryCommunities}, query)
	require.NoError(t, err)

	var communities []Community
	jsonErr := ModuleCodec.UnmarshalJSON(resBytes, &communities)
	assert.Nil(t, jsonErr)
	assert.Equal(t, communities[0].ID, "crypto")
	assert.Equal(t, communities[1].ID, "meme")
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
