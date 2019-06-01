package community

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()
	createdCommunity := keeper.NewCommunity(ctx, name, slug, description)

	params, err := json.Marshal(QueryCommunityParams{
		ID: 1,
	})
	assert.Nil(t, err)

	query := abci.RequestQuery{
		Path: "/custom/community/id",
		Data: params,
	}

	result, sdkErr := queryCommunity(ctx, query, keeper)
	assert.Nil(t, sdkErr)

	var returnedCommunity Community
	err = json.Unmarshal(result, &returnedCommunity)
	assert.Nil(t, err)
	assert.Equal(t, returnedCommunity.ID, createdCommunity.ID)
	assert.Equal(t, returnedCommunity.Name, createdCommunity.Name)
	assert.Equal(t, returnedCommunity.Slug, createdCommunity.Slug)
	assert.Equal(t, returnedCommunity.Description, createdCommunity.Description)
}

func TestQueryCommunity_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	params, err := json.Marshal(QueryCommunityParams{
		ID: 1,
	})
	require.Nil(t, err)

	query := abci.RequestQuery{
		Path: "/custom/community/id",
		Data: params,
	}

	_, sdkErr := queryCommunity(ctx, query, keeper)
	require.NotNil(t, sdkErr)
	require.Equal(t, ErrCommunityNotFound(1).Code(), sdkErr.Code(), "should get error")
}

func TestQueryCommunities_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()
	first := keeper.NewCommunity(ctx, name, slug, description)

	name2, slug2, description2 := getAnotherFakeCommunityParams()
	another := keeper.NewCommunity(ctx, name2, slug2, description2)

	result, sdkErr := queryCommunities(ctx, keeper)
	assert.Nil(t, sdkErr)

	var communities []Community
	err := json.Unmarshal(result, &communities)
	assert.Nil(t, err)
	assert.Equal(t, communities[0].ID, first.ID)
	assert.Equal(t, communities[1].ID, another.ID)
}
