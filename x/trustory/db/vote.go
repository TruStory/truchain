package db

import (
	"strconv"
	"strings"

	ts "github.com/TruStory/trucoin/x/trustory/types"
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
	err := k.cdc.UnmarshalBinary(storyVal, story)
	if err != nil {
		panic(err)
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
	voteVal, err := k.cdc.MarshalBinary(vote)
	if err != nil {
		panic(err)
	}
	voteStore.Set(voteKey, voteVal)

	// add vote id to story
	story.VoteIDs = append(story.VoteIDs, vote.ID)

	// create new story with vote
	newStory, err := k.cdc.MarshalBinary(*story)
	if err != nil {
		panic(err)
	}

	// replace old story with new one in story store
	storyStore.Set(storyKey, newStory)

	// add vote to vote list
	votes, err := k.GetActiveVotes(ctx, story.ID)
	if err != nil {
		panic(err)
	}
	votes = append(votes, vote.ID)
	err = k.SetActiveVotes(ctx, story.ID, votes)

	return vote.ID, nil
}

// ============================================================================

// GetActiveVotes gets all votes for the current round of a story
func (k TruKeeper) GetActiveVotes(ctx sdk.Context, storyID int64) ([]int64, sdk.Error) {
	store := ctx.KVStore(k.storyKey)
	key := generateVoteListKey(storyID)
	val := store.Get(key)
	if val == nil {
		return []int64{}, nil // FIXME: add error
	}
	votes := &[]int64{}
	err := k.cdc.UnmarshalBinary(val, votes)
	if err != nil {
		panic(err)
	}
	return *votes, nil
}

// SetActiveVotes sets all votes for the current round of a story
func (k TruKeeper) SetActiveVotes(ctx sdk.Context, storyID int64, votes []int64) sdk.Error {
	store := ctx.KVStore(k.storyKey)
	key := generateVoteListKey(storyID)
	value := k.cdc.MustMarshalBinary(votes)
	store.Set(key, value)

	return nil
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
	err := k.cdc.UnmarshalBinary(val, vote)
	if err != nil {
		panic(err)
	}
	return *vote, nil
}

// generateVoteListKey creates a key for a vote list of form "stories|ID|votes"
func generateVoteListKey(storyID int64) []byte {
	return []byte(strings.Join([]string{"stories", strconv.Itoa(int(storyID)), "votes"}, ""))
}
