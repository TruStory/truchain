package backing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processExpiredBackings(ctx)
	if err != nil {
		panic(err)
	}

	return sdk.NewTags()
}

// ============================================================================

// processExpiredBackings checks to see if backings have expired
func (k Keeper) processExpiredBackings(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", "backing")

	// find all expired backings
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

		if backing.IsExpired(ctx.BlockHeader().Time) {
			indicesToDelete = append(indicesToDelete, index)

			// distribute earnings from expired backings
			err := k.distributeEarnings(ctx, backing)
			if err != nil {
				panic(err)
			}
		}
		return false
	})

	for _, v := range indicesToDelete {
		backingList.Delete(v)
		msg := "Removed expired backing %d from backing list"
		logger.Info(fmt.Sprintf(msg, backingID))
	}

	return nil
}

// distributeEarnings adds coins from the backing to the user.
// Earnings is the original amount (principal) + interest.
func (k Keeper) distributeEarnings(ctx sdk.Context, backing Backing) sdk.Error {
	logger := ctx.Logger().With("module", "backing")

	// give the principal back to the user in category coins
	_, _, err := k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{backing.Amount()})
	if err != nil {
		return err
	}

	// give the interest earned to the user in category coins
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
