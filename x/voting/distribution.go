package voting

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	tokenVote "github.com/TruStory/truchain/x/vote"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) distributeRewards(
	ctx sdk.Context,
	pool sdk.Coin,
	votes poll,
	confirmed bool,
	categoryID int64) sdk.Error {

	logger := ctx.Logger().With("module", StoreKey)

	var winners []app.Voter
	if confirmed {
		winners = votes.trueVotes
	} else {
		winners = votes.falseVotes
	}

	stakerTotalAmount, stakerCount, voterCount, err := winnerInfo(winners)
	if err != nil {
		return err
	}

	var stakerPool sdk.Coin
	var voterPool sdk.Coin

	switch {
	case stakerCount > 0 && voterCount == 0:
		// staker pool is 100% of reward pool with no voters
		stakerPool = pool
	case stakerCount == 0 && voterCount > 0:
		// voter pool is 100% of reward pool with no stakers
		voterPool = pool
	default:
		// calculate reward pool for stakers (75% of pool)
		stakerPool = k.calculateStakerPool(ctx, pool)
		// calculate reward pool for voters (25% of pool)
		voterPool = k.calculateVoterPool(ctx, pool)
	}

	logger.Info(fmt.Sprintf("Staker reward pool: %v", stakerPool))
	logger.Info(fmt.Sprintf("Voter reward pool: %v", voterPool))

	// calculate voter reward amount
	voterRewardAmount := voterRewardAmount(voterPool, voterCount)

	for _, vote := range winners {
		logger.Info(fmt.Sprintf("Processing vote type: %T", vote))

		switch v := vote.(type) {
		case backing.Backing, challenge.Challenge:
			rewardCoin, err := k.rewardStaker(
				ctx, v, stakerTotalAmount, stakerPool, categoryID)
			if err != nil {
				return err
			}

			// remove reward amount from pool
			pool = pool.Minus(rewardCoin)

		case tokenVote.TokenVote:
			rewardCoin, err := k.rewardTokenVoter(ctx, v, voterRewardAmount, categoryID)
			if err != nil {
				return err
			}

			// remove reward amount from pool
			pool = pool.Minus(rewardCoin)
		}
	}

	logger.Info(fmt.Sprintf("Amount left in pool: %v", pool))

	err = checkForEmptyPool(pool, stakerCount+voterCount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) rewardStaker(
	ctx sdk.Context,
	staker app.Voter,
	stakerTotalAmount sdk.Int,
	stakerPool sdk.Coin,
	categoryID int64) (rewardCoin sdk.Coin, err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	// get back staked amount
	err = k.stakeKeeper.DistributePrincipalAndInterest(ctx, []app.Voter{staker}, categoryID)
	if err != nil {
		return
	}

	// calculate reward (X% of pool, in proportion to stake)
	rewardAmount := stakerRewardAmount(staker.Amount(), stakerTotalAmount, stakerPool)

	// calculate reward
	rewardCoin = sdk.NewCoin(app.StakeDenom, rewardAmount)

	// distribute reward in cred
	_, err = k.truBankKeeper.MintAndAddCoin(ctx, staker.Creator(), categoryID, rewardAmount)
	if err != nil {
		return
	}
	logger.Info(fmt.Sprintf("Distributed reward %s to staker", rewardCoin))

	return rewardCoin, nil
}

func (k Keeper) rewardTokenVoter(
	ctx sdk.Context,
	staker app.Voter,
	voterRewardAmount sdk.Int,
	categoryID int64) (rewardCoin sdk.Coin, err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	// get back original staked amount
	_, _, err = k.bankKeeper.AddCoins(ctx, staker.Creator(), sdk.Coins{staker.Amount()})
	if err != nil {
		return
	}
	logger.Info(fmt.Sprintf("Giving back original vote stake: %v", staker.Amount()))

	// calculate reward (1-X% of pool, in equal proportions)
	rewardCoin = sdk.NewCoin(app.StakeDenom, voterRewardAmount)

	// distribute reward in cred
	_, err = k.truBankKeeper.MintAndAddCoin(ctx, staker.Creator(), categoryID, voterRewardAmount)
	if err != nil {
		return
	}
	logger.Info(fmt.Sprintf("Distributed reward %s to voter", rewardCoin))

	return rewardCoin, nil
}

// winnerInfo returns data needed to calculate the reward pool
func winnerInfo(
	winners []app.Voter) (
	stakerTotalAmount sdk.Int,
	stakerCount int64,
	voterCount int64,
	err sdk.Error) {

	stakerTotalAmount = sdk.ZeroInt()

	for _, vote := range winners {
		switch v := vote.(type) {
		case backing.Backing, challenge.Challenge:
			stakerCount = stakerCount + 1
			stakerTotalAmount = stakerTotalAmount.Add(v.Amount().Amount)
		case tokenVote.TokenVote:
			voterCount = voterCount + 1
		default:
			return stakerTotalAmount, stakerCount, voterCount, ErrInvalidVote(v)
		}
	}

	return stakerTotalAmount, stakerCount, voterCount, nil
}
