package vote

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func rejectedPool(
	ctx sdk.Context, trueVotes []interface{}, falseVotes []interface{}, pool *sdk.Coin) (
	err sdk.Error) {

	// people who voted TRUE / lost the game
	for _, vote := range trueVotes {
		switch v := vote.(type) {

		case backing.Backing:
			// forfeit backing and inflationary rewards, add to pool
			*pool = (*pool).Plus(v.Amount).Plus(v.Interest)

		case app.Vote:
			// add vote fee to reward pool
			*pool = (*pool).Plus(v.Amount)

		default:
			if err = ErrVoteHandler(v); err != nil {
				return err
			}
		}
	}

	// people who voted FALSE / won the game
	for _, vote := range falseVotes {
		switch v := vote.(type) {

		case backing.Backing:
			// slash inflationary rewards and add to pool, bad boy
			*pool = (*pool).Plus(v.Interest)

		case challenge.Challenge:
			// do nothing
			// winning challengers keep their stake

		case app.Vote:
			// do nothing
			// winning voters keep their stake

		default:
			if err = ErrVoteHandler(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func distributeRewardsRejected(
	ctx sdk.Context, bankKeeper bank.Keeper, winners []interface{}, pool sdk.Coin) (err sdk.Error) {

	// load default parameters
	params := DefaultParams()

	// calculate reward pool for challengers (75% of pool)
	challengerPool := challengerPool(pool, params)

	// calculate reward pool for voters (25% of pool)
	voterPool := voterPool(pool, params)

	// count challengers and voters
	challengerCount, voterCount, err := count(winners)
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
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{v.Amount})

		case challenge.Challenge:
			// get back staked amount
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{v.Amount})
			if err != nil {
				return err
			}

			// get reward (X% of pool, in proportion to stake)
			rewardAmount := challengerRewardAmount(
				v.Amount, challengerCount, challengerPool)

			// mint coin and give money
			rewardCoin := sdk.NewCoin(pool.Denom, rewardAmount)
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{rewardCoin})

		case app.Vote:
			// get back original staked amount
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{v.Amount})
			if err != nil {
				return err
			}

			// get reward (1-X% of pool, in equal proportions)
			rewardCoin := sdk.NewCoin(pool.Denom, voterRewardAmount)
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{rewardCoin})

		default:
			err = ErrVoteHandler(v)
		}

		if err != nil {
			return err
		}
	}

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

// count challengers and voters
func count(
	winners []interface{}) (challengerCount int64, voterCount int64, err sdk.Error) {

	for _, vote := range winners {
		switch v := vote.(type) {
		case challenge.Challenge:
			challengerCount = challengerCount + 1
		case app.Vote:
			voterCount = voterCount + 1
		default:
			return 0, 0, ErrVoteHandler(v)
		}
	}

	return challengerCount, voterCount, nil
}

// amount / challengerCount * challengerPool
func challengerRewardAmount(
	amount sdk.Coin, challengerCount int64, challengerPool sdk.Coin) sdk.Int {

	return amount.Amount.
		Div(sdk.NewInt(challengerCount)).
		Mul(challengerPool.Amount)
}
