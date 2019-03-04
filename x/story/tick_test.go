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

func Test_processStoryQueue(t *testing.T) {
	ctx, storyKeeper := fakeStories()

	q := storyKeeper.pendingStoryQueue(ctx)
	assert.Equal(t, uint64(3), q.List.Len())

	err := storyKeeper.processStoryQueue(ctx, q)
	assert.Nil(t, err)

	story, _ := storyKeeper.Story(ctx, 5)
	assert.Equal(t, Pending, story.Status)

	// fake a future block time to expire story
	expiredTime := time.Now().Add(DefaultParams().ExpireDuration)
	ctx = ctx.WithBlockHeader(abci.Header{Time: expiredTime})

	err = storyKeeper.processStoryQueue(ctx, q)
	assert.Nil(t, err)

	story, _ = storyKeeper.Story(ctx, 3)
	assert.Equal(t, Expired, story.Status)
}
