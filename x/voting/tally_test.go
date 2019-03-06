package voting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// [CONFIRMED STORY] ==================================================
func TestTally(t *testing.T) {
	ctx, _, k := fakeConfirmedGame()

	storyID := int64(1)
	votes, _ := k.tally(ctx, storyID)

	assert.Equal(t, 7, len(votes.trueVotes))
	assert.Equal(t, 4, len(votes.falseVotes))
}

func TestConfirmStory(t *testing.T) {
	ctx, votes, k := fakeConfirmedGame()

	confirmed, _ := k.confirmStory(ctx, votes, "trudex", 1)
	assert.True(t, confirmed)
}

func TestWeightedVote(t *testing.T) {
	ctx, votes, k := fakeConfirmedGame()

	trueWeights, _ := k.weightedVote(ctx, votes.trueVotes, "trudex")
	falseWeights, _ := k.weightedVote(ctx, votes.falseVotes, "trudex")

	// 1 preethi added each due to cold-start
	assert.Equal(t, "7", trueWeights.String())

	// 1 preethi added each due to cold-start
	assert.Equal(t, "4", falseWeights.String())
}

// [REJECTED STORY] ==================================================

func TestTallyRejected(t *testing.T) {
	ctx, _, k := fakeRejectedGame()

	storyID := int64(1)
	votes, _ := k.tally(ctx, storyID)

	assert.Equal(t, 0, len(votes.trueVotes))
	assert.Equal(t, 1, len(votes.falseVotes))
}

func TestRejectedStory(t *testing.T) {
	ctx, votes, k := fakeRejectedGame()

	confirmed, _ := k.confirmStory(ctx, votes, "crypto", 1)
	assert.False(t, confirmed)
}

func TestWeightedVoteRejected(t *testing.T) {
	ctx, votes, k := fakeRejectedGame()

	trueWeights, _ := k.weightedVote(ctx, votes.trueVotes, "trudex")
	falseWeights, _ := k.weightedVote(ctx, votes.falseVotes, "trudex")

	// 1 preethi added each due to cold-start
	assert.Equal(t, "0", trueWeights.String())

	// 1 preethi added each due to cold-start
	assert.Equal(t, "1", falseWeights.String())
}
