package distribution

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_handleExpiredStoriesEmptyQueue(t *testing.T) {
	ctx, k, _, _, _, _ := mockDB()

	err := k.handleExpiredStories(ctx)
	assert.Nil(t, err)
}

func Test_handleExpiredStories(t *testing.T) {
	ctx, k, storyKeeper, backingKeeper, challengeKeeper, bankKeeper := mockDB()

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now()})
	storyID := createFakeStory(ctx, storyKeeper)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(100000))
	argument := "test argument right here"
	backer := fakeFundedCreator(ctx, bankKeeper)
	challenger := fakeFundedCreator(ctx, bankKeeper)
	duration := 5 * 24 * time.Hour

	_, err := backingKeeper.Create(
		ctx, storyID, amount, argument, backer, duration)
	assert.Nil(t, err)

	_, err = challengeKeeper.Create(
		ctx, storyID, amount, argument, challenger)
	assert.Nil(t, err)

	// fake expired story queue
	k.expiredStoryQueue(ctx).Push(storyID)

	err = k.handleExpiredStories(ctx)
	assert.Nil(t, err)

	// check distribution for backer
	coins := bankKeeper.GetCoins(ctx, backer)
	assert.Equal(t, "6670crypto,2000000000000trusteak", coins.String())

	// check balance for challenger
	coins = bankKeeper.GetCoins(ctx, challenger)
	assert.Equal(t, "2000000000000trusteak", coins.String())
}
