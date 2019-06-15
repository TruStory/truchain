package slashing

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerySlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	createdSlash, err := keeper.NewSlash(ctx, stakeID, creator)
	assert.Nil(t, err)

	params, jsonErr := json.Marshal(QuerySlashParams{
		ID: 1,
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: "/custom/slashing/id",
		Data: params,
	}

	result, sdkErr := querySlash(ctx, query, keeper)
	assert.Nil(t, sdkErr)

	var returnedSlash Slash
	jsonErr = json.Unmarshal(result, &returnedSlash)
	assert.Nil(t, jsonErr)
	assert.Equal(t, returnedSlash, createdSlash)
}

func TestQuerySlash_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	params, jsonErr := json.Marshal(QuerySlashParams{
		ID: 1,
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: "/custom/slashing/id",
		Data: params,
	}

	_, sdkErr := querySlash(ctx, query, keeper)
	assert.NotNil(t, sdkErr)
	assert.Equal(t, ErrSlashNotFound(1).Code(), sdkErr.Code())
}

func TestQuerySlashes_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	first, err := keeper.NewSlash(ctx, stakeID, creator)
	assert.Nil(t, err)

	stakeID2 := uint64(2)
	another, err := keeper.NewSlash(ctx, stakeID2, creator)
	assert.Nil(t, err)

	result, sdkErr := querySlashes(ctx, keeper)
	assert.Nil(t, sdkErr)

	var all []Slash
	jsonErr := json.Unmarshal(result, &all)
	assert.Nil(t, jsonErr)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], first)
	assert.Equal(t, all[1], another)
}
