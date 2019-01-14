package backing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	err := processBacking(ctx, k)
	if err != nil {
		panic(err)
	}

	return sdk.NewTags()
}

// ============================================================================

// processBacking checks each backing to see if it has expired. It calls itself
// recursively until all backings have been processed.
func processBacking(ctx sdk.Context, k Keeper) sdk.Error {
	// check if the backing queue is empty
	backing, err := k.QueueHead(ctx)
	if err != nil {
		if err.Code() == ErrQueueEmpty().Code() {
			return nil
		}
		return err
	}

	// check if backing has expired
	if ctx.BlockHeader().Time.Before(backing.Expires) {
		// no more expired backings left in queue
		// terminate recursion
		return nil
	}

	// remove expired backing from the queue
	if _, err = k.QueuePop(ctx); err != nil {
		return err
	}

	// distribute earnings to the backing creator
	if err = distributeEarnings(ctx, k, backing); err != nil {
		return err
	}

	// process next in queue
	return processBacking(ctx, k)
}

// distributeEarnings adds coins from the backing to the user.
// Earnings is the original amount (principal) + interest.
func distributeEarnings(ctx sdk.Context, k Keeper, backing Backing) sdk.Error {
	logger := ctx.Logger().With("module", "x/backing")

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
