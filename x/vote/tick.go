package vote

import (
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	store := ctx.KVStore(k.activeGamesQueueKey)
	q := queue.NewQueue(k.GetCodec(), store)

	err := checkGames(ctx, k, q)
	if err != nil {
		panic(err)
	}

	// TODO: maybe tags should return err?

	return sdk.NewTags()
}

// ============================================================================

// checkGames checks to see if a validation game has ended.
// It calls itself recursively until all games have been processed.
func checkGames(ctx sdk.Context, k Keeper, q queue.Queue) (err sdk.Error) {
	// 	// check the head of the queue
	// 	var challengeID int64
	// 	if err := q.Peek(&challengeID); err != nil {
	// 		return nil
	// 	}

	// 	// retrieve challenge from kvstore
	// 	game, err := k.challengeKeeper.Get(ctx, challengeID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// terminate recursion on finding the first unfinished challenge
	// 	if game.EndTime.After(ctx.BlockHeader().Time) {
	// 		return nil
	// 	}

	// 	// remove finished challenge from queue
	// 	q.Pop()

	// 	// tally backings, challenges, and votes
	// 	err = tally(ctx, k, game)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// TODO: distribute rewards

	// 	return checkGames(ctx, k, q)
	return
}

// func tally(
// 	ctx sdk.Context, keeper Keeper, game challenge.Game) (err sdk.Error) {

// 	// win is votes weighted by amount
// 	// noVotes = []interface{} (Backing, Challenge, Vote)
// 	// yesVotes = []interface{} (Backing, Challenge, Vote)

// 	// map ensures there can be no double counting
// 	var votes map[string]interface{}

// 	// tally backings
// 	err = tallyBackings(ctx, keeper.backingKeeper, game.StoryID, votes)
// 	if err != nil {
// 		return
// 	}

// 	// tally challenges
// 	err = tallyChallenges(keeper.challengeKeeper)
// 	if err != nil {
// 		return
// 	}

// 	// tally votes
// 	err = tallyVotes(keeper)
// 	if err != nil {
// 		return
// 	}

// 	// Vote type
// 	// Backing embeds Vote
// 	// Challenge embeds Vote

// 	for _, vote := range votes {
// 		switch vote.(type) {
// 		case backing.Backing:
// 			//
// 		case challenge.Challenge:
// 			//
// 		case app.Vote:
// 			//
// 		default:
// 			return sdk.ErrInternal("invalid type")
// 		}
// 	}

// 	return
// }
