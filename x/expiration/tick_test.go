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

	err := k.processStoryQueue(ctx)
	assert.Nil(t, err)
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
		ctx, storyID, amount, 0, argument, backer, false)
	assert.Nil(t, err)

	_, err = challengeKeeper.Create(
		ctx, storyID, amount, 0, argument, challenger, false)
	assert.Nil(t, err)

	// fake expired story queue
	k.storyQueue(ctx).Push(storyID)

	// fake future block time for expiration
	// expireTime := time.Now().Add(24 * time.Hour)
	expireTime := time.Now().Add(1 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: expireTime})

	err = k.processStoryQueue(ctx)
	assert.Nil(t, err)

	// check expiration for backer
	coins := bankKeeper.GetCoins(ctx, backer)
	assert.Equal(t, "1999999913340trusteak", coins.String())

	// check balance for challenger
	coins = bankKeeper.GetCoins(ctx, challenger)
	assert.Equal(t, "2000000313340trusteak", coins.String())
}
