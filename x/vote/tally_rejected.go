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

func rejectedPool(
	ctx sdk.Context, votes poll, pool *sdk.Coin) (err sdk.Error) {

	// people who voted TRUE / lost the game
	for _, vote := range votes.trueVotes {
		switch v := vote.(type) {

		case backing.Backing:
			// forfeit backing and inflationary rewards, add to pool
			// TODO [shanev]: do proper conversion when we know it, still 1:1
			interestInTrustake := sdk.NewCoin(params.StakeDenom, v.Interest.Amount)
			*pool = (*pool).Plus(v.Amount()).Plus(interestInTrustake)

		case TokenVote:
			// add vote fee to reward pool
			*pool = (*pool).Plus(v.Amount())

		default:
			if err = ErrInvalidVote(v); err != nil {
				return err
			}
		}
	}

	// people who voted FALSE / won the game
	for _, vote := range votes.falseVotes {
		switch v := vote.(type) {

		case backing.Backing:
			// slash inflationary rewards and add to pool, bad boy
			// TODO [shanev]: do proper conversion when we know it, still 1:1
			interestInTrustake := sdk.NewCoin(params.StakeDenom, v.Interest.Amount)
			*pool = (*pool).Plus(interestInTrustake)

		case challenge.Challenge:
			// do nothing
			// winning challengers keep their stake

		case TokenVote:
			// do nothing
			// winning voters keep their stake

		default:
			if err = ErrInvalidVote(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func distributeRewardsRejected(
	ctx sdk.Context,
	backingKeeper backing.WriteKeeper,
	bankKeeper bank.Keeper,
	votes poll,
	pool sdk.Coin,
	denom string) (err sdk.Error) {

	logger := ctx.Logger().With("module", "vote")

	// load default parameters
	defaults := DefaultParams()

	// get the total challenger stake amount and voter count
	challengerTotalAmount, challengerCount, voterCount, err :=
		winnerInfo(votes.falseVotes)
	if err != nil {
		return err
	}

	// challenger pool is 100% of reward pool when no voters
	challengerPool := pool
	voterPool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())

	if voterCount > 0 {
		// calculate reward pool for challengers (75% of pool)
		challengerPool = calculateChallengerPool(pool, defaults)
		logger.Info(fmt.Sprintf("Challenger reward pool: %v", challengerPool))

		// calculate reward pool for voters (25% of pool)
		voterPool = calculateVoterPool(pool, defaults)
		logger.Info(fmt.Sprintf("Voter reward pool: %v", voterPool))
	}

	// calculate voter reward amount
	voterRewardAmount := voterRewardAmount(voterPool, voterCount)

	// slash losers (true voters)
	for _, vote := range votes.trueVotes {
		logger.Info(fmt.Sprintf("Processing vote type: %T", vote))

		switch v := vote.(type) {
		case backing.Backing:
			// don't get anything back, too bad sucka!

		case challenge.Challenge:
			// challengers cannot vote true -- skip

		case TokenVote:
			// slashed -- get nothing back

		default:
			err = ErrInvalidVote(v)
		}

		if err != nil {
			return err
		}
	}
	// distribute reward to winners (false voters)
	for _, vote := range votes.falseVotes {
		logger.Info(fmt.Sprintf("Processing vote type: %T", vote))

		switch v := vote.(type) {
		case backing.Backing:
			// get back stake amount because we are nice
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Giving back original backing amount: %v", v.Amount()))

		case challenge.Challenge:
			// get back staked amount
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Giving back original challenge stake: %v", v.Amount()))

			// calculate reward (X% of pool, in proportion to stake)
			rewardAmount := challengerRewardAmount(
				v.Amount(), challengerTotalAmount, challengerPool)

			// calculate reward
			rewardCoin := sdk.NewCoin(pool.Denom, rewardAmount)

			// remove reward amount from pool
			pool = pool.Minus(rewardCoin)

			// distribute reward in cred
			cred := app.NewCategoryCoin(denom, rewardCoin)
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{cred})
			logger.Info(fmt.Sprintf("Distributed challenge reward of: %v", cred))

		case TokenVote:
			// get back original staked amount
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Giving back original vote stake: %v", v.Amount()))

			// calculate reward (1-X% of pool, in equal proportions)
			rewardCoin := sdk.NewCoin(pool.Denom, voterRewardAmount)

			// remove reward amount from pool
			pool = pool.Minus(rewardCoin)

			// distribute reward in cred
			cred := app.NewCategoryCoin(denom, rewardCoin)
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{cred})
			logger.Info(fmt.Sprintf("Distributed vote reward of: %v", cred))

		default:
			err = ErrInvalidVote(v)
		}

		if err != nil {
			return err
		}
	}

	logger.Info(fmt.Sprintf("Amount left in pool: %v", pool))

	err = checkForEmptyPoolRejected(pool, challengerCount, voterCount)
	if err != nil {
		return err
	}

	return nil
}

// calculate reward pool for challengers (75% of pool)
func calculateChallengerPool(pool sdk.Coin, params Params) sdk.Coin {

	challengerPoolShare := params.ChallengerRewardPoolShare

	challengerPoolAmount :=
		sdk.NewDecFromInt(pool.Amount).Mul(challengerPoolShare)

	return sdk.NewCoin(pool.Denom, challengerPoolAmount.TruncateInt())
}

// calculate reward pool for voters (25% of pool)
func calculateVoterPool(pool sdk.Coin, params Params) sdk.Coin {

	challengerPoolShare := params.ChallengerRewardPoolShare
	voterPoolShare := sdk.OneDec().Sub(challengerPoolShare)

	voterPoolAmount :=
		sdk.NewDecFromInt(pool.Amount).Mul(voterPoolShare)

	return sdk.NewCoin(pool.Denom, voterPoolAmount.TruncateInt())
}

// winnerInfo returns data needed to calculate the reward pool
func winnerInfo(
	winners []app.Voter) (
	challengerTotalAmount sdk.Int,
	challengerCount int64,
	voterCount int64,
	err sdk.Error) {

	challengerTotalAmount = sdk.ZeroInt()

	for _, vote := range winners {
		switch v := vote.(type) {
		case backing.Backing:
			// skip
		case challenge.Challenge:
			challengerCount = challengerCount + 1
			challengerTotalAmount = challengerTotalAmount.Add(v.Amount().Amount)
		case TokenVote:
			voterCount = voterCount + 1
		default:
			return challengerTotalAmount, challengerCount, voterCount, ErrInvalidVote(v)
		}
	}

	return challengerTotalAmount, challengerCount, voterCount, nil
}

// amount / challengerTotalAmount * challengerPool
func challengerRewardAmount(
	amount sdk.Coin, challengerTotalAmount sdk.Int, challengerPool sdk.Coin) sdk.Int {

	amountDec := sdk.NewDecFromInt(amount.Amount)

	rewardAmountDec := amountDec.
		QuoInt(challengerTotalAmount).
		MulInt(challengerPool.Amount)

	return rewardAmountDec.TruncateInt()
}

// accounts for leeway in reward pool due to division
func checkForEmptyPoolRejected(
	pool sdk.Coin, challengerCount int64, voterCount int64) sdk.Error {

	return checkForEmptyPoolConfirmed(pool, challengerCount+voterCount)
}
