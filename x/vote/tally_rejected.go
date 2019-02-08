package vote

import (
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
			*pool = (*pool).Plus(v.Amount()).Plus(v.Interest)

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
			*pool = (*pool).Plus(v.Interest)

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
	winners []app.Voter,
	pool sdk.Coin,
	denom string) (err sdk.Error) {

	// load default parameters
	params := DefaultParams()

	// calculate reward pool for challengers (75% of pool)
	challengerPool := challengerPool(pool, params)

	// calculate reward pool for voters (25% of pool)
	voterPool := voterPool(pool, params)

	// get the total challenger stake amount and voter count
	challengerTotalAmount, voterCount, err := winnerInfo(winners)
	if err != nil {
		return err
	}

	// calculate voter reward amount
	voterRewardAmount := voterRewardAmount(voterPool, voterCount)

	// distribute reward
	for _, vote := range winners {
		switch v := vote.(type) {

		case backing.Backing:
			// get back stake amount because we are nice
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})

			// remove backing from backing list
			err = backingKeeper.RemoveFromList(ctx, v.ID())

		case challenge.Challenge:
			// get back staked amount
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}

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

		case TokenVote:
			// get back original staked amount
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{v.Amount()})
			if err != nil {
				return err
			}

			// calculate reward (1-X% of pool, in equal proportions)
			rewardCoin := sdk.NewCoin(pool.Denom, voterRewardAmount)

			// remove reward amount from pool
			pool = pool.Minus(rewardCoin)

			// distribute reward in cred
			cred := app.NewCategoryCoin(denom, rewardCoin)
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator(), sdk.Coins{cred})

		default:
			err = ErrInvalidVote(v)
		}

		if err != nil {
			return err
		}
	}

	// TODO [shanev]: Remove after fixing https://github.com/TruStory/truchain/issues/314
	// err = checkForEmptyPool(pool)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// calculate reward pool for challengers (75% of pool)
func challengerPool(pool sdk.Coin, params Params) sdk.Coin {

	challengerPoolShare := params.ChallengerRewardPoolShare

	challengerPoolAmount :=
		sdk.NewDecFromInt(pool.Amount).Mul(challengerPoolShare)

	return sdk.NewCoin(pool.Denom, challengerPoolAmount.RoundInt())
}

// calculate reward pool for voters (25% of pool)
func voterPool(pool sdk.Coin, params Params) sdk.Coin {

	challengerPoolShare := params.ChallengerRewardPoolShare
	voterPoolShare := sdk.OneDec().Sub(challengerPoolShare)

	voterPoolAmount :=
		sdk.NewDecFromInt(pool.Amount).Mul(voterPoolShare)

	return sdk.NewCoin(pool.Denom, voterPoolAmount.RoundInt())
}

// winnerInfo returns data needed to calculate the reward pool
func winnerInfo(
	winners []app.Voter) (
	challengerTotalAmount sdk.Int, voterCount int64, err sdk.Error) {

	challengerTotalAmount = sdk.ZeroInt()

	for _, vote := range winners {
		switch v := vote.(type) {
		case backing.Backing:
			// skip
		case challenge.Challenge:
			challengerTotalAmount = challengerTotalAmount.Add(v.Amount().Amount)
		case TokenVote:
			voterCount = voterCount + 1
		default:
			return challengerTotalAmount, voterCount, ErrInvalidVote(v)
		}
	}

	return challengerTotalAmount, voterCount, nil
}

// amount / challengerTotalAmount * challengerPool
func challengerRewardAmount(
	amount sdk.Coin, challengerTotalAmount sdk.Int, challengerPool sdk.Coin) sdk.Int {

	amountDec := sdk.NewDecFromInt(amount.Amount)

	rewardAmountDec := amountDec.
		QuoInt(challengerTotalAmount).
		MulInt(challengerPool.Amount)

	return rewardAmountDec.RoundInt()
}
