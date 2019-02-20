package vote

import (
	"fmt"

	params "github.com/TruStory/truchain/parameters"
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// calculate reward pool for a confirmed story
func confirmedPool(
	ctx sdk.Context, falseVotes []app.Voter, pool *sdk.Coin) (err sdk.Error) {

	for _, vote := range falseVotes {
		switch v := vote.(type) {

		case backing.Backing:
			// slash inflationary rewards and add to pool
			// TODO [shanev]: do proper conversion when we know it, still 1:1
			interestInTrustake := sdk.NewCoin(params.StakeDenom, v.Interest.Amount)
			*pool = (*pool).Plus(interestInTrustake)

		case challenge.Challenge:
			// add challenge amount to reward pool
			*pool = (*pool).Plus(v.Amount())

		case TokenVote:
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

func distributeRewardsConfirmed(
	ctx sdk.Context,
	backingKeeper backing.WriteKeeper,
	bankKeeper bank.Keeper,
	votes poll,
	pool sdk.Coin,
	denom string) (err sdk.Error) {

	logger := ctx.Logger().With("module", "vote")

	// determine pool share per voter
	voterCount := int64(len(votes.trueVotes))
	voterRewardAmount := voterRewardAmount(pool, voterCount)
	rewardCoin := sdk.NewCoin(pool.Denom, voterRewardAmount)
	cred := app.NewCategoryCoin(denom, rewardCoin)

	logger.Info(fmt.Sprintf(
		"Token voter reward amount: %s", voterRewardAmount))

	// distribute reward to winners
	for _, vote := range votes.trueVotes {
		logger.Info(fmt.Sprintf("Processing winning vote type: %T", vote))

		switch v := vote.(type) {
		case backing.Backing:
			// distribute backing principal and interest
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf(
				"Giving back original backing principal: %v", v.Amount()))

			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Interest})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf(
				"Distributing backing interest: %v", v.Interest))

			pool = pool.Minus(rewardCoin)

			// distribute reward in cred
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{cred})
			logger.Info(fmt.Sprintf(
				"Distributed to backer a reward of: %v", cred))

		case TokenVote:
			// get back original staked amount in trustake
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf(
				"Giving back original vote amount: %v", v.Amount()))

			pool = pool.Minus(rewardCoin)

			// distribute reward in cred
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{cred})
			logger.Info(fmt.Sprintf(
				"Distributed to voter a reward of: %v", cred))

		default:
			if err = ErrInvalidVote(v); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}

	// slash losers
	for _, vote := range votes.falseVotes {
		logger.Info(fmt.Sprintf("Processing losing vote type: %T", vote))

		switch v := vote.(type) {
		// backer who changed their implicit TRUE vote to FALSE
		case backing.Backing:
			// return backing because we are nice people
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Giving back original backing amount: %v", v.Amount()))

		case challenge.Challenge:
			// do nothing
			// don't get their stake back

		case TokenVote:
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
