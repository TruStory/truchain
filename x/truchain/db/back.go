package db

import (
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewBacking adds a new vote to the vote store
func (k TruKeeper) NewBacking(
	ctx sdk.Context,
	story ts.Story,
	amount sdk.Coins,
	creator sdk.AccAddress,
	duration time.Duration,
) (int64, sdk.Error) {

	// create new backing type
	backing := ts.NewBacking(
		k.newID(ctx, k.backingKey),
		amount,
		time.Now().Add(duration),
		creator)

	// get handle for backing store
	store := ctx.KVStore(k.backingKey)

	// save it in the store
	store.Set(
		generateKey(k.backingKey.String(), backing.ID),
		k.cdc.MustMarshalBinary(backing))

	return backing.ID, nil
}

// GetVote gets a vote with the given id from the key-value store
// func (k TruKeeper) GetVote(ctx sdk.Context, voteID int64) (ts.Vote, sdk.Error) {
// 	store := ctx.KVStore(k.voteKey)
// 	key := generateKey(k.voteKey.String(), voteID)
// 	val := store.Get(key)
// 	if val == nil {
// 		return ts.Vote{}, ts.ErrVoteNotFound(voteID)
// 	}
// 	vote := &ts.Vote{}
// 	k.cdc.MustUnmarshalBinary(val, vote)

// 	return *vote, nil
// }

// ============================================================================

// GetBackers gets all votes for the current round of a story
// func (k TruKeeper) GetActiveVotes(ctx sdk.Context, storyID int64) []int64 {
// 	store := ctx.KVStore(k.storyKey)
// 	key := generateVoteListKey(storyID)
// 	val := store.Get(key)
// 	if val == nil {
// 		return []int64{}
// 	}
// 	votes := &[]int64{}
// 	k.cdc.MustUnmarshalBinary(val, votes)

// 	return *votes
// }

// SetBackers sets all votes for the current round of a story
// func (k TruKeeper) SetActiveVotes(ctx sdk.Context, storyID int64, votes []int64) {
// 	store := ctx.KVStore(k.storyKey)
// 	key := generateVoteListKey(storyID)
// 	value := k.cdc.MustMarshalBinary(votes)
// 	store.Set(key, value)
// }

// ============================================================================

// ============================================================================

// generateVoteListKey creates a key for a vote list of form "stories|ID|votes"
// func generateVoteListKey(storyID int64) []byte {
// 	return []byte(strings.Join([]string{"stories", strconv.Itoa(int(storyID)), "votes"}, ""))
// }
