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

func Test_processVotingStoryListNotMeetQuorum(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	err := k.processVotingStoryList(ctx)
	assert.Nil(t, err)
}

func Test_processVotingStoryListNotMeetVoteEndTime(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	k.votingStoryList(ctx).Push(int64(1))

	err := k.processVotingStoryList(ctx)
	assert.Nil(t, err)
}

func Test_processVotingStoryListVerifyStory(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	k.votingStoryList(ctx).Push(int64(1))

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)})

	err := k.processVotingStoryList(ctx)
	assert.Nil(t, err)
}

func Test_quorum(t *testing.T) {
	ctx, votes, k := fakeValidationGame()

	storyID := int64(1)
	totalBCV, _ := k.quorum(ctx, storyID)

	assert.Equal(t, len(votes.falseVotes)+len(votes.trueVotes), totalBCV)
}
