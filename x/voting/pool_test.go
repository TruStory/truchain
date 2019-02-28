package voting

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

// [CONFIRMED STORY] ==================================================
func TestConfirmedRewardPool(t *testing.T) {
	ctx, votes, k := fakeConfirmedGame()
	categoryID := int64(1)

	// fake future block time
	interestStopTime := ctx.BlockHeader().Time.Add(24 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: interestStopTime})

	pool, _ := k.rewardPool(ctx, votes, true, categoryID)
	assert.Equal(t, "5266999800000trusteak", pool.String())
}

func TestStakerRewardAmount(t *testing.T) {
	coin := stakerRewardAmount(
		sdk.NewCoin("trudex", sdk.NewInt(1000000000000)),
		sdk.NewInt(2000000000000),
		sdk.NewCoin("trudex", sdk.NewInt(5250000000000)))

	assert.Equal(t, "2625000000000", coin.String())
}

func TestStakerRewardAmount2(t *testing.T) {
	coin := stakerRewardAmount(
		sdk.NewCoin("trudex", sdk.NewInt(1500000000000)),
		sdk.NewInt(2000000000000),
		sdk.NewCoin("trudex", sdk.NewInt(5250000000000)))

	assert.Equal(t, "3937500000000", coin.String())
}

func TestStakerRewardAmount3(t *testing.T) {
	coin := stakerRewardAmount(
		sdk.NewCoin("trudex", sdk.NewInt(500000000000)),
		sdk.NewInt(2000000000000),
		sdk.NewCoin("trudex", sdk.NewInt(5250000000000)))

	assert.Equal(t, "1312500000000", coin.String())
}

func TestCheckForEmptyPool(t *testing.T) {
	pool, _ := sdk.ParseCoin("4trusteak")
	voterCount := int64(10)
	err := checkForEmptyPool(pool, voterCount)
	assert.Nil(t, err)
}

func TestCheckForEmptyPool2(t *testing.T) {
	pool, _ := sdk.ParseCoin("5trusteak")
	voterCount := int64(10)
	err := checkForEmptyPool(pool, voterCount)
	assert.Nil(t, err)
}

func TestCheckForEmptyPool3(t *testing.T) {
	pool, _ := sdk.ParseCoin("9trusteak")
	voterCount := int64(10)
	err := checkForEmptyPool(pool, voterCount)
	assert.Nil(t, err)
}

func Test_voterRewardAmount(t *testing.T) {
	pool, _ := sdk.ParseCoin("1trusteak")
	assert.Equal(t, sdk.NewInt(0), voterRewardAmount(pool, 0))
}

// [REJECTED STORY] ==================================================

func TestRejectedStoryRewardPool(t *testing.T) {
	ctx, votes, k := fakeRejectedGame()
	categoryID := int64(1)

	// fake future block time
	interestStopTime := ctx.BlockHeader().Time.Add(24 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: interestStopTime})

	pool, _ := k.rewardPool(ctx, votes, false, categoryID)
	assert.Equal(t, "0trusteak", pool.String())
}

func TestRejectedStakerPool(t *testing.T) {
	ctx, votes, k := fakeRejectedGame()
	categoryID := int64(1)

	// fake future block time
	interestStopTime := ctx.BlockHeader().Time.Add(24 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: interestStopTime})

	pool, _ := k.rewardPool(ctx, votes, false, categoryID)

	coin := k.calculateStakerPool(ctx, pool)
	assert.Equal(t, "0trusteak", coin.String())
}

func TestRejectedVoterPool(t *testing.T) {
	ctx, votes, k := fakeRejectedGame()
	categoryID := int64(1)

	// fake future block time
	interestStopTime := ctx.BlockHeader().Time.Add(24 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: interestStopTime})

	pool, _ := k.rewardPool(ctx, votes, false, categoryID)

	coin := k.calculateVoterPool(ctx, pool)
	assert.Equal(t, "0trusteak", coin.String())
}
