package vote

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// calculate reward pool for a confirmed story
func confirmedStoryRewardPool(
	ctx sdk.Context, bankKeeper bank.Keeper, falseVotes []interface{}) (pool sdk.Coin, err sdk.Error) {

	for _, vote := range falseVotes {
		switch v := vote.(type) {

		case backing.Backing:
			// slash inflationary rewards and add to pool
			pool = pool.Plus(v.Interest)

		case challenge.Challenge:
			// add challenge amount to reward pool
			pool = pool.Plus(v.Amount)

		case app.Vote:
			// add vote fee to reward pool
			pool = pool.Plus(v.Amount)

		default:
			err = ErrVoteHandler(v)
		}

		if err != nil {
			return pool, err
		}
	}

	return pool, nil
}

func distributeConfirmedStoryRewards(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	winners []interface{},
	losers []interface{},
	pool sdk.Coin) (err sdk.Error) {

	// determine pool share per voter
	voterRewardAmount := voterRewardAmount(pool, voterCount(winners))

	// distribute reward to winners
	for _, vote := range winners {
		switch v := vote.(type) {

		case backing.Backing:
			// keep backing as is

		case app.Vote:
			// get back original staked amount
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{v.Amount})

			// get money, an equal portion of the reward pool
			rewardCoin := sdk.NewCoin(pool.Denom, voterRewardAmount)
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{rewardCoin})

		default:
			err = ErrVoteHandler(v)
		}

		if err != nil {
			return err
		}
	}

	// slash losers
	for _, vote := range losers {
		switch v := vote.(type) {

		// backer who changed their implicit TRUE vote to FALSE
		case backing.Backing:
			// return backing because we are nice people
			_, _, err = bankKeeper.AddCoins(ctx, v.Creator, sdk.Coins{v.Amount})

		case challenge.Challenge:
			// do nothing
			// don't get their stake back

		case app.Vote:
			// do nothing
			// don't get their stake back

		default:
			err = ErrVoteHandler(v)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// count voters
func voterCount(winners []interface{}) (voterCount int64) {
	for _, voter := range winners {
		if _, ok := voter.(app.Vote); ok {
			voterCount = voterCount + 1
		}
	}

	return voterCount
}

func voterRewardAmount(pool sdk.Coin, voterCount int64) sdk.Int {

	return pool.Amount.Div(sdk.NewInt(voterCount))
}
