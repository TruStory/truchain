package category

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryCategories_ErrNotFound(t *testing.T) {
	ctx, k := mockDB()

	queryParams := QueryCategoryParams{
		ID: "1",
	}

	cdc := codec.New()

	bz, errRes := cdc.MarshalJSON(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/category/id",
		Data: bz,
	}

	_, err := queryCategoryByID(ctx, query, k)
	require.NotNil(t, err)
	require.Equal(t, ErrCategoryNotFound(1).Code(), err.Code(), "should get error")
}

func TestQueryCategoriesWithID(t *testing.T) {
	ctx, k := mockDB()

	createFakeCategory(ctx, k)

	queryParams := QueryCategoryParams{
		ID: "1",
	}

	cdc := codec.New()

	bz, errRes := cdc.MarshalJSON(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/category/id",
		Data: bz,
	}

	_, err := queryCategoryByID(ctx, query, k)
	require.Nil(t, err)
}
