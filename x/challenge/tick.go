package challenge

import (
	"fmt"

	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	err := k.filterPendingGameQueue(ctx, k.pendingGameQueue(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// filterPendingGameQueue checks to see if a pending voting game has expired
func (k Keeper) filterPendingGameQueue(ctx sdk.Context, pendingQueue queue.Queue) sdk.Error {
	logger := ctx.Logger().With("module", "challenge")

	if pendingQueue.IsEmpty() {
		return nil
	}

	var gameID int64
	if err := pendingQueue.Peek(&gameID); err != nil {
		panic(err)
	}

	game, err := k.gameKeeper.Game(ctx, gameID)
	if err != nil {
		return err
	}

	if game.IsExpired(ctx.BlockHeader().Time) {
		pendingQueue.Pop()
		msg := "Removed expired game %d from pending queue, len %d\n"
		logger.Info(fmt.Sprintf(msg, gameID, pendingQueue.List.Len()))

		k.returnFunds(ctx, gameID)

		// check next game in pending queue
		return k.filterPendingGameQueue(ctx, pendingQueue)
	}

	// done handling all expired games
	// exit recursion
	return nil
}

func (k Keeper) returnFunds(ctx sdk.Context, gameID int64) sdk.Error {
	logger := ctx.Logger().With("module", "challenge")

	// get challenges
	challenges, err := k.ChallengesByGameID(ctx, gameID)
	if err != nil {
		return err
	}

	// return funds
	for _, v := range challenges {
		msg := "Returned challenged amount back to %s for game %d."
		logger.Info(fmt.Sprintf(msg, v.Amount(), gameID))

		_, _, err = k.bankKeeper.AddCoins(
			ctx, v.Creator(), sdk.Coins{v.Amount()})
		if err != nil {
			return err
		}
	}

	return nil
}
