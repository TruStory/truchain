package db

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tendermint/tendermint/libs/log"
)

func TestVoteStory(t *testing.T) {
	ms, _, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	keeper := NewTruKeeper(storyKey, voteKey, bank.Keeper{}, cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	storyID := createFakeStory(ms, keeper)

	creator := sdk.AccAddress([]byte{3, 4})
	vote := true
	stake, _ := sdk.ParseCoins("10trustake")

	// test voting on a non-existant story
	_, err := keeper.VoteStory(ctx, int64(5), creator, vote, stake)
	assert.NotNil(t, err)

	// test voting on a story
	voteID, err := keeper.VoteStory(ctx, storyID, creator, vote, stake)
	assert.Nil(t, err)
	assert.Equal(t, voteID, int64(0), "Vote ID does not match")

	// test getting a non-existant vote
	_, err = keeper.GetVote(ctx, int64(5))
	assert.NotNil(t, err)

	// test getting vote and comparing fields
	savedVote, err := keeper.GetVote(ctx, voteID)
	assert.Nil(t, err)
	assert.Equal(t, savedVote.Vote, true, "Vote choice  does not match")

	assert.Equal(t, savedVote.Amount.AmountOf("trustake"), sdk.NewInt(10), "Vote amount does not match")
}
