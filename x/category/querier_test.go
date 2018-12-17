package category

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryCategories_ErrNotFound(t *testing.T) {
	ctx, k := mockDB()

	queryParams := QueryCategoryByIDParams{
		ID: 1,
	}

	bz, errRes := json.Marshal(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/categories/id",
		Data: bz,
	}

	_, err := queryCategoryByID(ctx, query, k)
	require.NotNil(t, err)
	require.Equal(t, ErrCategoryNotFound(1).Code(), err.Code(), "should get error")
}

func TestQueryCategoriesWithID(t *testing.T) {
	ctx, k := mockDB()

	createFakeCategory(ctx, k)

	queryParams := QueryCategoryByIDParams{
		ID: 1,
	}

	bz, errRes := json.Marshal(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/categories/id",
		Data: bz,
	}

	_, err := queryCategoryByID(ctx, query, k)

	require.Nil(t, err)
}
