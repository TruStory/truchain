package distribution

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.handleExpiredStories(ctx)
	if err != nil {
		panic(err)
	}
	return sdk.EmptyTags()
}

// recursively process expired stories
func (k Keeper) handleExpiredStories(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", "distribution")

	expiredStoryQueue := k.expiredStoryQueue(ctx)

	if expiredStoryQueue.IsEmpty() {
		// done processing all expired stories
		// terminate
		return nil
	}

	var storyID int64
	if err := expiredStoryQueue.Peek(&storyID); err != nil {
		panic(err)
	}
	logger.Info(fmt.Sprintf("Handling expired story id: %d", storyID))

	expiredStoryQueue.Pop()

	err := k.distributeEarningsToBackers(ctx, storyID)
	if err != nil {
		return err
	}

	err = k.returnFundsToChallengers(ctx, storyID)
	if err != nil {
		return err
	}

	// handle next expired story
	return k.handleExpiredStories(ctx)
}

func (k Keeper) distributeEarningsToBackers(ctx sdk.Context, storyID int64) sdk.Error {
	logger := ctx.Logger().With("module", "distribution")

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	backings, err := k.backingKeeper.BackingsByStoryID(ctx, storyID)
	if err != nil {
		return err
	}

	for _, backing := range backings {
		// give the principal back to the user (in trustake)
		_, _, err := k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{backing.Amount()})
		if err != nil {
			return err
		}

		// give the interest earned to the user (in cred)
		period := story.VotingEndTime.Sub(story.Timestamp.CreatedTime)
		denom, err := k.storyKeeper.CategoryDenom(ctx, storyID)
		if err != nil {
			return err
		}
		interest := getInterest(
			backing.Amount(), period, period, denom, DefaultParams())
		_, _, err = k.bankKeeper.AddCoins(ctx, backing.Creator(), sdk.Coins{interest})
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf(
			"Distributed earnings of %s with interest of %s to %s",
			backing.Amount().String(),
			backing.Interest.String(),
			backing.Creator().String()))
	}

	return nil
}

func (k Keeper) returnFundsToChallengers(ctx sdk.Context, storyID int64) sdk.Error {
	logger := ctx.Logger().With("module", "distribution")

	// get challenges
	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, storyID)
	if err != nil {
		return err
	}

	// return funds
	for _, v := range challenges {
		msg := "Returned challenged amount %s back to %s for story %d."
		logger.Info(fmt.Sprintf(msg, v.Amount(), v.Creator(), storyID))

		_, _, err = k.bankKeeper.AddCoins(
			ctx, v.Creator(), sdk.Coins{v.Amount()})
		if err != nil {
			return err
		}
	}

	return nil
}

func getInterest(
	amount sdk.Coin,
	period time.Duration,
	maxPeriod time.Duration,
	credDenom string,
	params Params) sdk.Coin {

	// TODO: keep track of total supply
	// https://github.com/TruStory/truchain/issues/22

	totalSupply := sdk.NewDec(1000000000000000)

	// inputs
	maxAmount := totalSupply
	amountWeight := params.AmountWeight
	periodWeight := params.PeriodWeight
	maxInterestRate := params.MaxInterestRate

	// type cast values to unitless decimals for math operations
	periodDec := sdk.NewDec(int64(period))
	maxPeriodDec := sdk.NewDec(int64(maxPeriod))
	amountDec := sdk.NewDecFromInt(amount.Amount)

	// normalize amount and period to 0 - 1
	normalizedAmount := amountDec.Quo(maxAmount)
	normalizedPeriod := periodDec.Quo(maxPeriodDec)

	// apply weights to normalized amount and period
	weightedAmount := normalizedAmount.Mul(amountWeight)
	weightedPeriod := normalizedPeriod.Mul(periodWeight)

	// calculate interest
	interestRate := maxInterestRate.Mul(weightedAmount.Add(weightedPeriod))
	// convert rate to a value
	if interestRate.LT(params.MinInterestRate) {
		interestRate = params.MinInterestRate
	}
	interest := amountDec.Mul(interestRate)

	// return cred coin with rounded interest
	cred := sdk.NewCoin(credDenom, interest.RoundInt())

	return cred
}
