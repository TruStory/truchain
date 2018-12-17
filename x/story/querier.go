package story

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath                = "stories"
	QueryStoryByID           = "id"
	QueryStoriesByCategoryID = "category"
	QueryStories             = "all"
)

// QueryCategoryStoriesParams are params for stories by category queries
type QueryCategoryStoriesParams struct {
	CategoryID int64
}

// QueryStoryByIDParams are params for getting a story
type QueryStoryByIDParams struct {
	ID int64
}

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryStoryByID:
			return queryStoryByID(ctx, req, k)
		case QueryStoriesByCategoryID:
			return queryStoriesByCategoryID(ctx, req, k)
		case QueryStories:
			return queryStories(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown truchain query endpoint")
		}
	}
}

// ============================================================================

func queryStoryByID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	params := QueryStoryByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	story, err := k.Story(ctx, params.ID)
	if err != nil {
		return
	}

	return app.MustMarshal(story), nil
}

func queryStoriesByCategoryID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	params := QueryCategoryStoriesParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	stories, err := k.StoriesByCategoryID(ctx, params.CategoryID)
	if err != nil {
		return
	}

	return app.MustMarshal(stories), nil
}

func queryStories(ctx sdk.Context, _ abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	stories := k.Stories(ctx)

	return app.MustMarshal(stories), nil
}
