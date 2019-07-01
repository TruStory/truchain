package staking

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerier_EmptyTopArgument(t *testing.T) {
	ctx, k, _, _, _ := mockDB()

	querier := NewQuerier(k)
	queryParams := QueryClaimTopArgumentParams{
		ClaimID: 1,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryClaimTopArgument}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryClaimTopArgument}, query)
	assert.NoError(t, err)
	argument := Argument{}
	jsonErr := k.codec.UnmarshalJSON(bz, &argument)
	assert.NoError(t, jsonErr)
	assert.Equal(t, uint64(0), argument.ID)

}
