package voting

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestEndBlock(t *testing.T) {
	ctx, _, k := fakeConfirmedGame()

	tags := k.EndBlock(ctx)
	assert.Equal(t, sdk.Tags{}, tags)
}

func Test_processVotingStoryListNotMeetVoteEndTime(t *testing.T) {
	ctx, _, k := fakeConfirmedGame()

	k.challengedStoryQueue(ctx).Push(int64(1))

	err := k.processChallengeStoryQueue(ctx)
	assert.Nil(t, err)
}

func Test_processVotingStoryListVerifyStory(t *testing.T) {
	ctx, _, k := fakeConfirmedGame()

	k.challengedStoryQueue(ctx).Push(int64(1))

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)})

	err := k.processChallengeStoryQueue(ctx)
	assert.Nil(t, err)
}
