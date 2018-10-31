package story

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryStories_ErrNotFound(t *testing.T) {
	ctx, k, _ := mockDB()

	queryParams := QueryCategoryStoriesParams{
		CategoryID: "1",
	}

	cdc := codec.New()

	bz, errRes := cdc.MarshalJSON(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/category/stories",
		Data: bz,
	}

	_, err := queryStoriesWithCategory(ctx, query, k)
	require.NotNil(t, err)
	require.Equal(t, ErrStoriesWithCategoryNotFound(1).Code(), err.Code(), "should get error")
}

func TestQueryStoriesWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	createFakeStory(ctx, sk, ck)

	queryParams := QueryCategoryStoriesParams{
		CategoryID: "1",
	}

	cdc := codec.New()

	bz, errRes := cdc.MarshalJSON(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/category/stories",
		Data: bz,
	}

	_, err := queryStoriesWithCategory(ctx, query, sk)
	require.Nil(t, err)
}

func TestQueryChallengedStoriesWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	sk.StartChallenge(ctx, storyID)

	queryParams := QueryCategoryStoriesParams{
		CategoryID: "1",
	}

	cdc := codec.New()

	bz, errRes := cdc.MarshalJSON(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/category/stories",
		Data: bz,
	}

	_, err := queryChallengedStoriesWithCategory(ctx, query, sk)
	require.Nil(t, err)
}

func TestQueryStoryFeed(t *testing.T) {
	ctx, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	sk.StartChallenge(ctx, storyID)

	queryParams := QueryCategoryStoriesParams{
		CategoryID: "1",
	}

	cdc := codec.New()

	bz, errRes := cdc.MarshalJSON(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/category/stories",
		Data: bz,
	}

	_, err := queryStoryFeed(ctx, query, sk)
	require.Nil(t, err)
}
