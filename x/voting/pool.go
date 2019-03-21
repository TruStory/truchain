package voting

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/stake"
	tokenVote "github.com/TruStory/truchain/x/vote"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) rewardPool(
	ctx sdk.Context,
	votes poll,
	confirmed bool,
	categoryID int64) (pool sdk.Coin, err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	pool = sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())

	var losers []stake.Voter
	if confirmed {
		losers = votes.falseVotes
	} else {
		losers = votes.trueVotes
	}

	for _, vote := range losers {
		switch v := vote.(type) {

		case backing.Backing, challenge.Challenge:

			// period := ctx.BlockHeader().Time.Sub(v.Timestamp().CreatedTime)
			// interest := k.stakeKeeper.Interest(ctx, v.Amount(), period)
			// interestCoin := sdk.NewCoin(app.StakeDenom, interest)

			// // add stake + interest to reward pool
			// pool = pool.Plus(v.Amount()).Plus(interestCoin)

			// logger.Info(fmt.Sprintf(
			// 	"Added stake %s with %s interest to pool for staker",
			// 	v.Amount(), interestCoin))

		case tokenVote.TokenVote:

			// add vote fee to reward pool
			pool = pool.Plus(v.Amount())

			logger.Info(fmt.Sprintf(
				"Added stake %s to pool for voter", v.Amount()))
		}
	}

	return pool, nil
}

// calculate reward pool for stakers (75% of pool)
func (k Keeper) calculateStakerPool(ctx sdk.Context, pool sdk.Coin) sdk.Coin {
	stakerPoolAmount := sdk.NewDecFromInt(pool.Amount).
		Mul(k.stakerRewardPoolShare(ctx))

	return sdk.NewCoin(pool.Denom, stakerPoolAmount.TruncateInt())
}

// calculate reward pool for voters (25% of pool)
func (k Keeper) calculateVoterPool(ctx sdk.Context, pool sdk.Coin) sdk.Coin {
	voterPoolShare := sdk.OneDec().Sub(k.stakerRewardPoolShare(ctx))
	voterPoolAmount := sdk.NewDecFromInt(pool.Amount).Mul(voterPoolShare)

	return sdk.NewCoin(pool.Denom, voterPoolAmount.TruncateInt())
}

// amount / stakerTotalAmount * stakerPool
func stakerRewardAmount(
	amount sdk.Coin, stakerTotalAmount sdk.Int, stakerPool sdk.Coin) sdk.Int {

	return sdk.NewDecFromInt(amount.Amount).
		QuoInt(stakerTotalAmount).
		MulInt(stakerPool.Amount).
		TruncateInt()
}

// accounts for leeway in reward pool due to division
// 169 pool / 10 voters = 16.9 = 16 per voter, 9 left in pool
// pool must be < voter counter
func checkForEmptyPool(pool sdk.Coin, voterCount int64) sdk.Error {
	if pool.Amount.GT(sdk.NewInt(voterCount)) {
		return ErrNonEmptyRewardPool(pool)
	}

	return nil
}

func voterRewardAmount(pool sdk.Coin, voterCount int64) sdk.Int {
	// check for no token voters
	// prevent division by zero errors
	if voterCount == 0 {
		return sdk.NewInt(0)
	}

	poolDec := sdk.NewDecFromInt(pool.Amount)
	voterCountInt := sdk.NewInt(voterCount)

	return poolDec.
		QuoInt(voterCountInt).
		TruncateInt()
}
