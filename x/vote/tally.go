package vote

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/game"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// tally votes and distribute rewards
func processGame(ctx sdk.Context, k Keeper, game game.Game) sdk.Error {
	// tally backings, challenges, and votes
	votes, err := tally(ctx, k, game)
	if err != nil {
		return err
	}

	// check if story was confirmed
	confirmed, err := confirmStory(ctx, k.accountKeeper, votes)
	if err != nil {
		return err
	}

	// calculate reward pool
	rewardPool, err := rewardPool(ctx, votes, confirmed)
	if err != nil {
		return err
	}

	// distribute rewards
	err = distributeRewards(
		ctx, k.bankKeeper, rewardPool, votes, confirmed)
	if err != nil {
		return err
	}

	// update story state
	err = k.storyKeeper.EndGame(ctx, game.StoryID, confirmed)
	if err != nil {
		return err
	}

	return nil
}

// tally backings, challenges, and token votes into two true and false slices
func tally(ctx sdk.Context, k Keeper, game game.Game) (votes poll, err sdk.Error) {

	// tally backings
	trueBackings, falseBackings, err := k.backingKeeper.Tally(ctx, game.StoryID)
	if err != nil {
		return
	}
	for _, v := range trueBackings {
		votes.trueVotes = append(votes.trueVotes, v)
	}
	for _, v := range falseBackings {
		votes.falseVotes = append(votes.falseVotes, v)
	}

	// tally challenges
	falseChallenges, err := k.challengeKeeper.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	for _, v := range falseChallenges {
		votes.falseVotes = append(votes.falseVotes, v)
	}

	// tally token votes
	trueTokenVotes, falseTokenVotes, err := k.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	for _, v := range trueTokenVotes {
		votes.trueVotes = append(votes.trueVotes, v)
	}
	for _, v := range falseTokenVotes {
		votes.falseVotes = append(votes.falseVotes, v)
	}

	return votes, nil
}

func rewardPool(
	ctx sdk.Context, votes poll, confirmed bool) (pool sdk.Coin, err sdk.Error) {

	// initialize an empty reward pool, false votes will always exist
	// because challengers with implicit false votes will always exist
	v, ok := votes.falseVotes[0].(app.Voter)
	if !ok {
		return pool, ErrInvalidVote(v, "Initializing reward pool")
	}
	pool = sdk.NewCoin(v.Amount().Denom, sdk.ZeroInt())

	if confirmed {
		err = confirmedPool(ctx, votes.falseVotes, &pool)
	} else {
		err = rejectedPool(ctx, votes, &pool)
	}
	if err != nil {
		return pool, err
	}

	return pool, nil
}

func distributeRewards(
	ctx sdk.Context, bankKeeper bank.Keeper, rewardPool sdk.Coin, votes poll, confirmed bool) (
	err sdk.Error) {

	logger := ctx.Logger().With("module", "vote")

	if confirmed {
		err = distributeRewardsConfirmed(
			ctx, bankKeeper, votes, rewardPool)
	} else {
		err = distributeRewardsRejected(
			ctx, bankKeeper, votes.falseVotes, rewardPool)
	}
	if err != nil {
		return
	}

	logger.Info("Distributed reward pool of " + rewardPool.String())

	return
}

// determine if a story is confirmed or rejected
func confirmStory(
	ctx sdk.Context, accountKeeper auth.AccountKeeper, votes poll) (
	confirmed bool, err sdk.Error) {

	// calculate weighted true votes
	trueWeight, err := weightedVote(ctx, accountKeeper, votes.trueVotes)
	if err != nil {
		return confirmed, err
	}

	// calculate weighted false votes
	falseWeight, err := weightedVote(ctx, accountKeeper, votes.falseVotes)
	if err != nil {
		return confirmed, err
	}

	// calculate what percent of the total weight is true votes
	totalWeight := trueWeight.Add(falseWeight)
	trueWeightDec := sdk.NewDecFromInt(trueWeight)
	truePercentOfTotal := trueWeightDec.QuoInt(totalWeight)

	// majority weight wins
	if truePercentOfTotal.GTE(DefaultParams().MajorityPercent) {
		// story confirmed
		return true, nil
	}

	// story rejected
	return false, nil
}

// calculate weighted vote based on user's total category coin balance
func weightedVote(
	ctx sdk.Context, accountKeeper auth.AccountKeeper, votes []app.Voter) (
	weightedAmount sdk.Int, err sdk.Error) {

	weightedAmount = sdk.ZeroInt()

	// iterate through all types of votes
	for _, vote := range votes {
		v, ok := vote.(app.Voter)
		if !ok {
			return weightedAmount, ErrInvalidVote(v)
		}

		// get the user account
		user := accountKeeper.GetAccount(ctx, v.Creator())

		// get user's total amount of category coins for story category
		categoryCoins := user.GetCoins().AmountOf(v.Amount().Denom)

		// total up the amount of category coins all voters have
		weightedAmount = weightedAmount.Add(categoryCoins)
	}

	return weightedAmount, nil
}

// Make sure reward pool is empty (<= 1 coin)
// Makes up for rounding error during division
func checkForEmptyPool(pool sdk.Coin) sdk.Error {
	oneCoin := sdk.NewCoin(pool.Denom, sdk.OneInt())
	if !(pool.IsLT(oneCoin) || pool.IsEqual(oneCoin)) {
		return ErrNonEmptyRewardPool(pool)
	}

	return nil
}
