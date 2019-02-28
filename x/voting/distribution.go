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

	// staker pool is 100% of reward pool when no voters
	stakerPool := pool
	voterPool := sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())

	if voterCount > 0 {
		// calculate reward pool for stakers (75% of pool)
		stakerPool = k.calculateStakerPool(ctx, pool)
		logger.Info(fmt.Sprintf("Staker reward pool: %v", stakerPool))
		fmt.Printf("Staker reward pool: %v\n", stakerPool)

		// calculate reward pool for voters (25% of pool)
		voterPool = k.calculateVoterPool(ctx, pool)
		logger.Info(fmt.Sprintf("Voter reward pool: %v", voterPool))
		fmt.Printf("Voter reward pool: %v\n", voterPool)
	}

	// calculate voter reward amount
	voterRewardAmount := voterRewardAmount(voterPool, voterCount)

	for _, vote := range winners {
		logger.Info(fmt.Sprintf("Processing vote type: %T", vote))
		fmt.Printf("Processing vote type: %T\n", vote)

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
	fmt.Printf("Amount left in pool: %v\n", pool)

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
	fmt.Printf("Distributed reward %s to staker\n", rewardCoin)

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
	fmt.Printf("Giving back original vote stake: %v\n", staker.Amount())

	// calculate reward (1-X% of pool, in equal proportions)
	rewardCoin = sdk.NewCoin(app.StakeDenom, voterRewardAmount)

	// distribute reward in cred
	_, err = k.truBankKeeper.MintAndAddCoin(ctx, staker.Creator(), categoryID, voterRewardAmount)
	if err != nil {
		return
	}
	logger.Info(fmt.Sprintf("Distributed reward %s to voter", rewardCoin))
	fmt.Printf("Distributed reward %s to voter\n", rewardCoin)

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
