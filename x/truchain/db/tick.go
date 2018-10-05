package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewResponseEndBlock is called at the end of every block, processes timing
// related events, and returns a ResponseEndBlock.
func (k TruKeeper) NewResponseEndBlock(ctx sdk.Context) abci.ResponseEndBlock {
	err := processBacking(ctx, k)
	if err != nil {
		panic(err)
	}

	return abci.ResponseEndBlock{}
}

// ============================================================================

// processBacking checks each backing to see if it has expired. It calls itself
// recursively until all backings have been processed.
func processBacking(ctx sdk.Context, k TruKeeper) sdk.Error {

	// check if the backing queue is empty
	backing, err := k.BackingQueueHead(ctx)
	if err != nil {
		if err.Code() == ts.ErrBackingQueueEmpty().Code() {
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
	if _, err = k.BackingQueuePop(ctx); err != nil {
		return err
	}

	// distribute earnings to the backing creator
	distributeEarnings(ctx, k, backing)

	// process next in queue
	return processBacking(ctx, k)
}

// distributeEarnings adds coins from the backing to the user.
// Earnings is the original amount (principal) + interest.
func distributeEarnings(ctx sdk.Context, k TruKeeper, backing ts.Backing) sdk.Error {

	// give the principal back to the user in category coins
	_, _, err := k.ck.AddCoins(ctx, backing.User, sdk.Coins{backing.Principal})
	if err != nil {
		return err
	}

	// give the interest earned to the user in category coins
	_, _, err = k.ck.AddCoins(ctx, backing.User, sdk.Coins{backing.Interest})
	if err != nil {
		return err
	}

	return nil
}
