package voting

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestEndBlock(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	tags := k.EndBlock(ctx)
	assert.Equal(t, sdk.Tags{}, tags)
}

func Test_processVotingStoryListNotMeetVoteEndTime(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	k.votingStoryQueue(ctx).Push(int64(1))

	err := k.processVotingStoryQueue(ctx)
	assert.Nil(t, err)
}

func Test_processVotingStoryListVerifyStory(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	k.votingStoryQueue(ctx).Push(int64(1))

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)})

	err := k.processVotingStoryQueue(ctx)
	assert.Nil(t, err)
}
