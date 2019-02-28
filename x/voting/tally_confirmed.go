package voting

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	tokenVote "github.com/TruStory/truchain/x/vote"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// calculate reward pool for a confirmed story
func (k Keeper) confirmedPool(
	ctx sdk.Context,
	falseVotes []app.Voter,
	pool *sdk.Coin,
	categoryID int64) (err sdk.Error) {

	for _, vote := range falseVotes {
		switch v := vote.(type) {

		case backing.Backing:
			// slash inflationary rewards and add to pool
			period := ctx.BlockHeader().Time.Sub(v.Timestamp().CreatedTime)
			interest := k.stakeKeeper.Interest(ctx, v.Amount(), categoryID, period)
			interestCoin := sdk.NewCoin(app.StakeDenom, interest)
			*pool = (*pool).Plus(interestCoin)

		case challenge.Challenge:
			// add challenge amount to reward pool
			*pool = (*pool).Plus(v.Amount())

		case tokenVote.TokenVote:
			// add vote fee to reward pool
			*pool = (*pool).Plus(v.Amount())

		default:
			if err = ErrInvalidVote(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) distributeRewardsConfirmed(
	ctx sdk.Context,
	votes poll,
	pool sdk.Coin,
	categoryID int64) (err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	// determine pool share per voter
	voterCount := int64(len(votes.trueVotes))
	voterRewardAmount := voterRewardAmount(pool, voterCount)
	rewardCoin := sdk.NewCoin(pool.Denom, voterRewardAmount)

	logger.Info(fmt.Sprintf("Token voter reward amount: %s", voterRewardAmount))

	// distribute reward to winners
	for _, vote := range votes.trueVotes {
		logger.Info(fmt.Sprintf("Processing winning vote type: %T", vote))

		switch v := vote.(type) {
		case backing.Backing:
			err := k.stakeKeeper.DistributePrincipalAndInterest(ctx, []app.Voter{v}, categoryID)
			if err != nil {
				return err
			}

			pool = pool.Minus(rewardCoin)

			_, err = k.truBankKeeper.MintAndAddCoin(ctx, v.Creator(), categoryID, voterRewardAmount)
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Distributed reward %s to backer", rewardCoin))

		case tokenVote.TokenVote:
			_, _, err = k.bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Giving back original vote amount: %v", v.Amount()))

			pool = pool.Minus(rewardCoin)

			_, err = k.truBankKeeper.MintAndAddCoin(ctx, v.Creator(), categoryID, voterRewardAmount)
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Distributed reward %s to voter", rewardCoin))

		default:
			if err = ErrInvalidVote(v); err != nil {
				return err
			}
		}
	}

	// slash losers
	for _, vote := range votes.falseVotes {
		logger.Info(fmt.Sprintf("Processing losing vote type: %T", vote))

		switch v := vote.(type) {
		// backer who changed their implicit TRUE vote to FALSE
		case backing.Backing:
			// return backing because we are nice people
			_, _, err = k.bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Giving back original backing amount: %v", v.Amount()))

		case challenge.Challenge:
			// do nothing
			// don't get their stake back

		case tokenVote.TokenVote:
			// do nothing
			// don't get their stake back

		default:
			err = ErrInvalidVote(v)
		}

		if err != nil {
			return err
		}
	}

	logger.Info(fmt.Sprintf("Amount left in pool: %v", pool))

	err = checkForEmptyPoolConfirmed(pool, voterCount)
	if err != nil {
		return err
	}

	return nil
}

// accounts for leeway in reward pool due to division
// 169 pool / 10 voters = 16.9 = 16 per voter, 9 left in pool
// pool must be < voter counter
func checkForEmptyPoolConfirmed(pool sdk.Coin, voterCount int64) sdk.Error {
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
