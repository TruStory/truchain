package game

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_checkStories(t *testing.T) {
	ctx, gameKeeper, _, _, _, _ := mockDB()

	err := gameKeeper.checkStories(ctx)
	assert.Nil(t, err)
}

func Test_checkStoriesMeetsQuorumAndThreshold(t *testing.T) {
	ctx, gameKeeper, storyKeeper, backingKeeper, challengeKeeper, bankKeeper := mockDB()

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now()})
	storyID := createFakeStory(ctx, storyKeeper)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(100000))
	argument := "test argument right here"
	backer1 := fakeFundedCreator(ctx, bankKeeper)
	backer2 := fakeFundedCreator(ctx, bankKeeper)
	challenger1 := fakeFundedCreator(ctx, bankKeeper)
	challenger2 := fakeFundedCreator(ctx, bankKeeper)
	duration := 5 * 24 * time.Hour

	_, err := backingKeeper.Create(
		ctx, storyID, amount, argument, backer1, duration)
	assert.Nil(t, err)

	_, err = backingKeeper.Create(
		ctx, storyID, amount, argument, backer2, duration)
	assert.Nil(t, err)

	_, err = challengeKeeper.Create(
		ctx, storyID, amount, argument, challenger1)
	assert.Nil(t, err)

	_, err = challengeKeeper.Create(
		ctx, storyID, amount, argument, challenger2)
	assert.Nil(t, err)

	err = gameKeeper.checkStories(ctx)
	assert.Nil(t, err)

	story, _ := storyKeeper.Story(ctx, storyID)
	assert.Equal(t, story.State.String(), "Voting")
}
