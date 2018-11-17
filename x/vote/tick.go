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

	return sdk.NewTags()
}

// ============================================================================

// checkGames checks to see if a validation game has ended.
// It calls itself recursively until all games have been processed.
func checkGames(ctx sdk.Context, k Keeper, q queue.Queue) (err sdk.Error) {
	// check the head of the queue
	var gameID int64
	if err := q.Peek(&gameID); err != nil {
		return nil
	}

	// retrieve the game
	game, err := k.gameKeeper.Get(ctx, gameID)
	if err != nil {
		return err
	}

	// terminate recursion on finding the first non-ended game
	if game.Ended(ctx.BlockHeader().Time) {
		return nil
	}

	// remove ended game from queue
	q.Pop()

	// process ended game
	err = processGame(ctx, k, game)
	if err != nil {
		return err
	}

	// check next game
	return checkGames(ctx, k, q)
}
