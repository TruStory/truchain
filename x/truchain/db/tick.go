package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewResponseEndBlock checks stories and generates a ResponseEndBlock.
// It is called at the end of every block, and processes any timing-related
// acitivities within the app.
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
		// done processing queue, return with no error
		if err.Code() == ts.CodeBackingQueueEmpty {
			return nil
		}
		return err
	}

	// process next backing if this one hasn't expired
	if ctx.BlockHeader().Time.Before(backing.Expires) {
		return processBacking(ctx, k)
	}

	// remove expired backing from the queue
	k.BackingQueuePop(ctx)

	// distribute earnings to the backing creator
	distributeEarnings(ctx, k, backing)

	// done processing queue, return with no error
	return nil
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
