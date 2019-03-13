package story

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, storyKeeper := fakeStories()

	tags := storyKeeper.EndBlock(ctx)
	assert.Equal(t, sdk.Tags{}, tags)
}

func Test_processStoryList(t *testing.T) {
	ctx, storyKeeper := fakeStories()

	l := storyKeeper.pendingStoryList(ctx)
	assert.Equal(t, uint64(3), l.Len())

	err := storyKeeper.processPendingStoryList(ctx, l)
	assert.Nil(t, err)

	story, _ := storyKeeper.Story(ctx, 5)
	assert.Equal(t, Pending, story.Status)

	// fake a future block time to expire story
	expiredTime := time.Now().Add(DefaultParams().ExpireDuration)
	ctx = ctx.WithBlockHeader(abci.Header{Time: expiredTime})

	err = storyKeeper.processPendingStoryList(ctx, l)
	assert.Nil(t, err)

	story, _ = storyKeeper.Story(ctx, 3)
	assert.Equal(t, Expired, story.Status)
}

func Test_processStoryList_All(t *testing.T) {
	ctx, storyKeeper := fakeStories()
	l := storyKeeper.pendingStoryList(ctx)
	assert.Equal(t, uint64(3), l.Len())

	story2, err := storyKeeper.Story(ctx, 2)
	assert.NoError(t, err)
	story2.Status = Challenged
	storyKeeper.UpdateStory(ctx, story2)
	err = storyKeeper.processPendingStoryList(ctx, l)
	assert.Nil(t, err)

	var currentStoryID int64
	l.Iterate(&currentStoryID, func(index uint64) bool {
		assert.NotEqual(t, int64(2), currentStoryID, "storyId 2 should have been processed")
		return false
	})

	// fake a future block time to expire story
	expiredTime := time.Now().Add(DefaultParams().ExpireDuration)
	story3, err := storyKeeper.Story(ctx, 3)
	assert.NoError(t, err)
	storyKeeper.UpdateStory(ctx, story3)
	ctx = ctx.WithBlockHeader(abci.Header{Time: expiredTime})

	err = storyKeeper.processPendingStoryList(ctx, l)
	assert.Nil(t, err)

	l.Iterate(&currentStoryID, func(index uint64) bool {
		assert.NotEqual(t, int64(3), currentStoryID, "storyId 3 should have been expired")
		return false
	})

	story3, _ = storyKeeper.Story(ctx, 3)
	assert.Equal(t, Expired, story3.Status)

}
