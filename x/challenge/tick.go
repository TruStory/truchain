package challenge

import (
	"fmt"

	list "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.filterExpiredGames(ctx, k.pendingGameList(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// filterExpiredGames checks to see if a pending game has expired
func (k Keeper) filterExpiredGames(ctx sdk.Context, pendingList list.List) sdk.Error {
	// logger := ctx.Logger().With("module", "challenge")

	// find all expired games
	// var gameID int64
	// var indicesToDelete []uint64
	// pendingList.Iterate(&gameID, func(index uint64) bool {
	// 	var tempGameID int64
	// 	err := pendingList.Get(index, &tempGameID)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	game, err := k.gameKeeper.Game(ctx, tempGameID)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if game.IsExpired(ctx.BlockHeader().Time) {
	// 		indicesToDelete = append(indicesToDelete, index)
	// 	}

	// 	return false
	// })

	// // delete expired games and return challenged funds
	// for _, v := range indicesToDelete {
	// 	pendingList.Delete(v)

	// 	msg := "Removed expired game %d from pending game list"
	// 	logger.Info(fmt.Sprintf(msg, gameID))

	// 	err := k.returnFunds(ctx, gameID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (k Keeper) returnFunds(ctx sdk.Context, gameID int64) sdk.Error {
	logger := ctx.Logger().With("module", "challenge")

	// get challenges
	challenges, err := k.ChallengesByStoryID(ctx, gameID)
	if err != nil {
		return err
	}

	// return funds
	for _, v := range challenges {
		msg := "Returned challenged amount %s back to %s for game %d."
		logger.Info(fmt.Sprintf(msg, v.Amount(), v.Creator(), gameID))

		_, _, err = k.bankKeeper.AddCoins(
			ctx, v.Creator(), sdk.Coins{v.Amount()})
		if err != nil {
			return err
		}
	}

	return nil
}
