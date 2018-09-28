package db

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func TestVoteStory(t *testing.T) {
	ctx, ms, am, k := mockDB()

	// create and fund voter's account
	denom := "memecoin"
	amount := sdk.NewInt(5)
	coins := sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(5))}
	creatorAddr := createFundedAccount(ctx, am, coins)

	// fake params for voting on story
	storyID := createFakeStory(ms, k)
	vote := true

	// test voting on a non-existant story
	_, err := k.VoteStory(ctx, int64(5), creatorAddr, vote, coins)
	assert.NotNil(t, err)

	// test voting on a story
	voteID, err := k.VoteStory(ctx, storyID, creatorAddr, vote, coins)
	fmt.Print(err)
	assert.Nil(t, err)
	assert.Equal(t, voteID, int64(0), "Vote ID does not match")

	// test getting a non-existant vote
	_, err = k.GetVote(ctx, int64(5))
	assert.NotNil(t, err)

	// test getting vote and comparing fields
	savedVote, err := k.GetVote(ctx, voteID)
	assert.NotNil(t, err)
	assert.Equal(t, savedVote.Vote, true, "Vote choice  does not match")
	assert.Equal(t, savedVote.Amount.AmountOf(denom), amount, "Vote amount does not match")
}
