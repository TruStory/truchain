package expiration

import (
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_handleExpiredStoriesEmptyQueue(t *testing.T) {
	ctx, k, _, _, _, _ := mockDB()
	completedStories, err := k.processStoryQueue(ctx, make([]app.CompletedStory, 0))
	assert.Nil(t, err)
	assert.Len(t, completedStories, 0)
}

func Test_handleExpiredStories(t *testing.T) {
	ctx, k, storyKeeper, backingKeeper, challengeKeeper, bankKeeper := mockDB()

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now()})
	storyID := createFakeStory(ctx, storyKeeper)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	argument := "test argument right here"
	backer := fakeFundedCreator(ctx, bankKeeper)
	challenger := fakeFundedCreator(ctx, bankKeeper)

	_, err := backingKeeper.Create(
		ctx, storyID, amount, 0, argument, backer)
	assert.Nil(t, err)

	_, err = challengeKeeper.Create(
		ctx, storyID, amount, 0, argument, challenger)
	assert.Nil(t, err)

	// fake future block time for expiration
	expireTime := time.Now().Add(100 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: expireTime})

	completedStories, err := k.processStoryQueue(ctx, make([]app.CompletedStory, 0))
	assert.Nil(t, err)
	assert.Len(t, completedStories, 1)

	assert.Equal(t, completedStories[0].ID, storyID)
	assert.Equal(t, completedStories[0].Creator, sdk.AccAddress([]byte{1, 2}))
	assert.Equal(t, completedStories[0].Challengers, []app.Staker{app.Staker{Address: challenger, Amount: amount}})
	assert.Equal(t, completedStories[0].Backers, []app.Staker{app.Staker{Address: backer, Amount: amount}})
	assert.Equal(t, completedStories[0].StakeDistributionResults.TotalAmount, amount.Add(amount))
	assert.Equal(t, completedStories[0].StakeDistributionResults.Type, app.DistributionMajorityNotReached)

	// check expiration for backer
	coins := bankKeeper.GetCoins(ctx, backer)
	expectedCoin := sdk.NewCoin("trusteak", sdk.NewInt(20*app.Shanev))
	assert.True(t, coins.IsAllGT(sdk.Coins{expectedCoin}))

	// check balance for challenger
	coins = bankKeeper.GetCoins(ctx, challenger)
	assert.True(t, coins.IsAllGT(sdk.Coins{expectedCoin}))
}
