package db

import (
	"strconv"
	"strings"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VoteStory saves a vote to a story
func (k TruKeeper) VoteStory(ctx sdk.Context, storyID int64, creator sdk.AccAddress, choice bool, amount sdk.Coins) (int64, sdk.Error) {
	// access story DB
	storyStore := ctx.KVStore(k.storyKey)
	storyKey := generateKey(k.storyKey.String(), storyID)
	storyVal := storyStore.Get(storyKey)

	if storyVal == nil {
		return -1, ts.ErrStoryNotFound(storyID)
	}

	// get existing story
	story := &ts.Story{}
	k.cdc.MustUnmarshalBinary(storyVal, story)

	// temporarily moves funds from voter to an escrow account until
	// the voting period is over and funds are distributed
	_, err := k.ck.SendCoins(ctx, creator, story.Escrow, amount)
	if err != nil {
		return -1, err
	}

	// create new vote struct
	vote := ts.Vote{
		ID:           k.newID(ctx, k.voteKey),
		StoryID:      story.ID,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		Round:        story.Round + 1,
		Amount:       amount,
		Vote:         choice,
	}

	// store vote in vote store
	voteStore := ctx.KVStore(k.voteKey)
	voteKey := generateKey(k.voteKey.String(), vote.ID)
	voteVal := k.cdc.MustMarshalBinary(vote)
	voteStore.Set(voteKey, voteVal)

	// add vote id to story
	story.VoteIDs = append(story.VoteIDs, vote.ID)

	// create new story with vote
	newStory := k.cdc.MustMarshalBinary(*story)

	// replace old story with new one in story store
	storyStore.Set(storyKey, newStory)

	// add vote to vote list
	votes := k.GetActiveVotes(ctx, story.ID)
	votes = append(votes, vote.ID)
	k.SetActiveVotes(ctx, story.ID, votes)

	return vote.ID, nil
}

// ============================================================================

// GetActiveVotes gets all votes for the current round of a story
func (k TruKeeper) GetActiveVotes(ctx sdk.Context, storyID int64) []int64 {
	store := ctx.KVStore(k.storyKey)
	key := generateVoteListKey(storyID)
	val := store.Get(key)
	if val == nil {
		return []int64{}
	}
	votes := &[]int64{}
	k.cdc.MustUnmarshalBinary(val, votes)

	return *votes
}

// SetActiveVotes sets all votes for the current round of a story
func (k TruKeeper) SetActiveVotes(ctx sdk.Context, storyID int64, votes []int64) {
	store := ctx.KVStore(k.storyKey)
	key := generateVoteListKey(storyID)
	value := k.cdc.MustMarshalBinary(votes)
	store.Set(key, value)
}

// ============================================================================

// GetVote gets a vote with the given id from the key-value store
func (k TruKeeper) GetVote(ctx sdk.Context, voteID int64) (ts.Vote, sdk.Error) {
	store := ctx.KVStore(k.voteKey)
	key := generateKey(k.voteKey.String(), voteID)
	val := store.Get(key)
	if val == nil {
		return ts.Vote{}, ts.ErrVoteNotFound(voteID)
	}
	vote := &ts.Vote{}
	k.cdc.MustUnmarshalBinary(val, vote)

	return *vote, nil
}

// ============================================================================

// generateVoteListKey creates a key for a vote list of form "stories|ID|votes"
func generateVoteListKey(storyID int64) []byte {
	return []byte(strings.Join([]string{"stories", strconv.Itoa(int(storyID)), "votes"}, ""))
}
