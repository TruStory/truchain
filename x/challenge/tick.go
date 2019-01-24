package challenge

import (
	"fmt"

	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	pendingGameQueue := k.pendingGameQueue(ctx)
	err := k.checkPendingQueue(ctx, pendingGameQueue)
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// checkPendingQueue checks to see if a pending voting game has expired
func (k Keeper) checkPendingQueue(ctx sdk.Context, pendingQueue queue.Queue) sdk.Error {
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
		fmt.Printf("Removed expired game %d from pending queue\n", gameID)

		k.returnFunds(ctx, gameID)

		// check next game in pending queue
		return k.checkPendingQueue(ctx, pendingQueue)
	}

	// done handling all expired games
	// exit recursion
	return nil
}

func (k Keeper) returnFunds(ctx sdk.Context, gameID int64) sdk.Error {
	logger := ctx.Logger().With("module", "x/game")

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
