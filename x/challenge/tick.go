package challenge

import (
	store "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	unexpiredChallengeQueue := store.NewQueue(k.GetCodec(), k.GetStore(ctx))
	err := checkExpiredChallenges(ctx, k, unexpiredChallengeQueue)
	if err != nil {
		panic(err)
	}

	return sdk.NewTags()
}

// ============================================================================

// checkExpiredChallenges checks each challenge to see if it has expired.
// It calls itself recursively until all challenges have been processed.
func checkExpiredChallenges(ctx sdk.Context, k Keeper, q store.Queue) sdk.Error {
	// check the head of the queue
	var challengeID int64
	if err := q.Peek(&challengeID); err != nil {
		return nil
	}

	// retrieve challenge from kvstore
	challenge, err := k.Get(ctx, challengeID)
	if err != nil {
		return err
	}

	// all remaining challenges expire at a later time
	if challenge.ExpiresTime.After(ctx.BlockHeader().Time) {
		// terminate recursion
		return nil
	}

	// remove expired challenge from queue
	q.Pop()

	// return funds and delete challenge if it hasn't started
	if !challenge.Started {
		err = returnFunds(ctx, k, challenge)
		if err != nil {
			return err
		}
		k.delete(ctx, challengeID)

		// remove challenge association from story
		story, err := k.storyKeeper.GetStory(ctx, challenge.StoryID)
		if err != nil {
			return err
		}
		story.ChallengeID = 0
		k.storyKeeper.UpdateStory(ctx, story)
	}

	return checkExpiredChallenges(ctx, k, q)
}

// returnFunds iterates through the challenger keyspace and returns funds
func returnFunds(ctx sdk.Context, k Keeper, challenge Challenge) sdk.Error {
	store := k.GetStore(ctx)

	// builds prefix of from "challenges:id:5:userAddr:"
	prefix := k.getChallengerPrefix(challenge.ID)

	// iterates through keyspace to find all challengers on a challenge
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var challenger ChallengerInfo
		bz := iter.Value()
		if bz == nil {
			return ErrNotFoundChallenger(challenge.ID)
		}
		k.GetCodec().MustUnmarshalBinary(bz, &challenger)

		// return funds
		_, _, err := k.bankKeeper.AddCoins(
			ctx,
			challenger.User,
			sdk.Coins{challenger.Amount})
		if err != nil {
			return err
		}
	}

	return nil
}
