package slashing

import (
	"strings"
	"testing"

	"github.com/TruStory/truchain/x/staking"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerySlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	createdSlash, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NoError(t, err)

	params := keeper.codec.MustMarshalJSON(QuerySlashParams{
		ID: 1,
	})

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QuerySlash}, "/"),
		Data: params,
	}

	result, sdkErr := querySlash(ctx, query, keeper)
	assert.NoError(t, sdkErr)

	var returnedSlash Slash
	jsonErr := keeper.codec.UnmarshalJSON(result, &returnedSlash)
	assert.NoError(t, jsonErr)
	assert.Equal(t, returnedSlash, createdSlash)
}

func TestQuerySlash_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	params := keeper.codec.MustMarshalJSON(QuerySlashParams{
		ID: 1,
	})
	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QuerySlash}, "/"),
		Data: params,
	}

	_, sdkErr := querySlash(ctx, query, keeper)
	assert.NotNil(t, sdkErr)
	assert.Equal(t, ErrSlashNotFound(1).Code(), sdkErr.Code())
}

func TestQuerySlashes_Success(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[0]
	_, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeBacking)
	assert.NoError(t, err)

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	first, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)

	creator2 := keeper.GetParams(ctx).SlashAdmins[1]
	another, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator2)
	assert.Nil(t, err)

	result, sdkErr := querySlashes(ctx, keeper)
	assert.Nil(t, sdkErr)

	var all []Slash
	jsonErr := keeper.codec.UnmarshalJSON(result, &all)
	assert.NoError(t, jsonErr)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], first)
	assert.Equal(t, all[1], another)
}
