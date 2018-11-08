package game

import (
	store "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	// TODO [Shane] add back expired challenges checker
	// https://github.com/TruStory/truchain/issues/44

	return sdk.EmptyTags()
}

// ============================================================================

// checkExpiredGames checks each challenge to see if it has expired.
// It calls itself recursively until all challenges have been processed.
func checkExpiredGames(ctx sdk.Context, k Keeper, q store.Queue) sdk.Error {
	// check the head of the queue
	var gameID int64
	if err := q.Peek(&gameID); err != nil {
		return nil
	}

	// retrieve game from kvstore
	game, err := k.Get(ctx, gameID)
	if err != nil {
		return err
	}

	// all remaining challenges expire at a later time
	if game.ExpiresTime.After(ctx.BlockHeader().Time) {
		// terminate recursion
		return nil
	}

	// remove expired challenge from queue
	q.Pop()

	// return funds and delete challenge if it hasn't started
	// if !game.Started {
	// 	if err = returnFunds(ctx, k, game); err != nil {
	// 		return err
	// 	}
	// 	if err = k.delete(ctx, gameID); err != nil {
	// 		return err
	// 	}

	// 	// TODO [Shane]: also delete challengers here
	// 	// see https://github.com/TruStory/truchain/issues/54

	// 	// remove challenge association from story
	// 	story, err := k.storyKeeper.GetStory(ctx, game.StoryID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	story.GameID = 0
	// 	k.storyKeeper.UpdateStory(ctx, story)
	// }

	return checkExpiredGames(ctx, k, q)
}

// returnFunds iterates through the challenger keyspace and returns funds
// func returnFunds(ctx sdk.Context, k Keeper, game Game) sdk.Error {
// 	store := k.GetStore(ctx)

// 	// builds prefix of from "challenges:id:5:userAddr:"
// 	prefix := k.getChallengerPrefix(challenge.ID)

// 	// iterates through keyspace to find all challengers on a challenge
// 	iter := sdk.KVStorePrefixIterator(store, prefix)
// 	defer iter.Close()
// 	for ; iter.Valid(); iter.Next() {
// 		var challenger Challenger
// 		bz := iter.Value()
// 		if bz == nil {
// 			return ErrNotFoundChallenger(challenge.ID)
// 		}
// 		k.GetCodec().MustUnmarshalBinary(bz, &challenger)

// 		// return funds
// 		_, _, err := k.bankKeeper.AddCoins(
// 			ctx,
// 			challenger.Creator,
// 			sdk.Coins{challenger.Amount})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
