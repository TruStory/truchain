package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================
// Vote operations

// NewVote adds a new vote to the vote store
func (k TruKeeper) NewVote(
	ctx sdk.Context,
	story ts.Story,
	amount sdk.Coins,
	creator sdk.AccAddress,
	choice bool,
) (int64, sdk.Error) {

	// create new vote
	vote := ts.NewVote(k.newID(ctx, k.voteKey), story.ID, amount, ctx.BlockHeight(), creator, story.Round+1, choice)

	// save it in the store
	voteStore := ctx.KVStore(k.voteKey)
	voteKey := generateKey(k.voteKey.String(), vote.ID)
	voteVal := k.cdc.MustMarshalBinary(vote)
	voteStore.Set(voteKey, voteVal)

	return vote.ID, nil
}

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
