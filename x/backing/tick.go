package backing

import (
	"fmt"

	list "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processMaturedBackings(ctx)
	if err != nil {
		panic(err)
	}

	return sdk.NewTags()
}

// ============================================================================

// processMaturedBackings checks to see if backings have matured
func (k Keeper) processMaturedBackings(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", "backing")

	// find all matured backings
	var backingID int64
	var indicesToDelete []uint64
	backingList := k.backingList(ctx)
	backingList.Iterate(&backingID, func(index uint64) bool {
		var tempBackingID int64
		err := backingList.Get(index, &tempBackingID)
		if err != nil {
			panic(err)
		}

		backing, err := k.Backing(ctx, tempBackingID)
		if err != nil {
			panic(err)
		}

		if k.isGameInSession(ctx, backing.StoryID()) {
			// skip maturing it
			// process next one
			return false
		}

		if backing.HasMatured(ctx.BlockHeader().Time) {
			indicesToDelete = append(indicesToDelete, index)

			// distribute earnings from matured backings
			err := k.distributeEarnings(ctx, backing)
			if err != nil {
				panic(err)
			}
		}
		return false
	})

	for _, v := range indicesToDelete {
		backingList.Delete(v)
		msg := "Removed matured backing %d from backing list"
		logger.Info(fmt.Sprintf(msg, backingID))
	}

	return nil
}

func (k Keeper) isGameInSession(ctx sdk.Context, storyID int64) bool {
	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		panic(err)
	}

	gameID := story.GameID

	// check if game is going on...
	if gameID > 0 {
		gameFoundInPendingGameList := k.isGameInList(k.pendingGameList(ctx), gameID)
		gameFoundInGameQueue := k.isGameInList(k.gameQueue(ctx).List, gameID)
		if gameFoundInPendingGameList || gameFoundInGameQueue {
			// backing is in challenged or voting phase
			return true
		}
	}

	return false
}

func (k Keeper) isGameInList(gameList list.List, gameID int64) bool {
	var found bool
	var ID int64
	gameList.Iterate(&ID, func(uint64) bool {
		if ID == gameID {
			found = true
			return true
		}
		return false
	})

	return found
}

// distributeEarnings adds coins from the backing to the user.
// Earnings is the original amount (principal) + interest.
func (k Keeper) distributeEarnings(ctx sdk.Context, backing Backing) sdk.Error {
	logger := ctx.Logger().With("module", "backing")

	// give the principal back to the user (in trustake)
	_, _, err := k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{backing.Amount()})
	if err != nil {
		return err
	}

	// give the interest earned to the user (in cred)
	_, _, err = k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{backing.Interest})
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf(
		"Distributed earnings of %s with interest of %s to %s",
		backing.Amount().String(),
		backing.Interest.String(),
		backing.Creator().String()))

	return nil
}
