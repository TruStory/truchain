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
	trueVotes, falseVotes, err := tally(ctx, k, game)
	if err != nil {
		return err
	}

	// check if story was confirmed
	confirmed, err := confirmStory(ctx, k.accountKeeper, trueVotes, falseVotes)
	if err != nil {
		return err
	}

	// calculate reward pool
	rewardPool, err := rewardPool(ctx, trueVotes, falseVotes, confirmed)
	if err != nil {
		return err
	}

	// distribute rewards
	err = distributeRewards(
		ctx, k.bankKeeper, rewardPool, trueVotes, falseVotes, confirmed)
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
func tally(
	ctx sdk.Context, k Keeper, game game.Game) (
	trueVotes []interface{}, falseVotes []interface{}, err sdk.Error) {

	// tally backings
	trueBackings, falseBackings, err := k.backingKeeper.Tally(ctx, game.StoryID)
	if err != nil {
		return
	}
	for _, v := range trueBackings {
		trueVotes = append(trueVotes, v)
	}
	for _, v := range falseBackings {
		falseVotes = append(falseVotes, v)
	}

	// tally challenges
	falseChallenges, err := k.challengeKeeper.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	for _, v := range falseChallenges {
		falseVotes = append(falseVotes, v)
	}

	// tally token votes
	trueTokenVotes, falseTokenVotes, err := k.Tally(ctx, game.ID)
	if err != nil {
		return
	}
	for _, v := range trueTokenVotes {
		trueVotes = append(trueVotes, v)
	}
	for _, v := range falseTokenVotes {
		falseVotes = append(falseVotes, v)
	}

	return trueVotes, falseVotes, nil
}

func rewardPool(
	ctx sdk.Context, trueVotes []interface{}, falseVotes []interface{}, confirmed bool) (
	pool sdk.Coin, err sdk.Error) {

	// initialize an empty reward pool, false votes will always exist
	// because challengers with implicit false votes will always exist
	v, ok := falseVotes[0].(app.Voter)
	if !ok {
		return pool, ErrInvalidVote(v, "Initializing reward pool")
	}
	pool = sdk.NewCoin(v.AmountDenom(), sdk.ZeroInt())

	if confirmed {
		err = confirmedPool(ctx, falseVotes, &pool)
	} else {
		err = rejectedPool(ctx, trueVotes, falseVotes, &pool)
	}
	if err != nil {
		return pool, err
	}

	return pool, nil
}

func distributeRewards(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	rewardPool sdk.Coin,
	trueVotes []interface{},
	falseVotes []interface{},
	confirmed bool) (err sdk.Error) {

	if confirmed {
		err = distributeRewardsConfirmed(
			ctx, bankKeeper, trueVotes, falseVotes, rewardPool)
	} else {
		err = distributeRewardsRejected(
			ctx, bankKeeper, falseVotes, rewardPool)
	}
	if err != nil {
		return
	}

	return
}

// determine if a story is confirmed or rejected
func confirmStory(
	ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	trueVotes []interface{},
	falseVotes []interface{}) (confirmed bool, err sdk.Error) {

	// calculate weighted votes
	trueWeight, err := weightedVote(ctx, accountKeeper, trueVotes)
	if err != nil {
		return confirmed, err
	}

	falseWeight, err := weightedVote(ctx, accountKeeper, falseVotes)
	if err != nil {
		return confirmed, err
	}

	totalWeight := trueWeight.Add(falseWeight)
	trueWeightDec := sdk.NewDecFromInt(trueWeight)
	truePercentOfTotal := trueWeightDec.QuoInt(totalWeight)

	// supermajority wins
	if truePercentOfTotal.GT(DefaultParams().SupermajorityPercent) {
		// story confirmed
		return true, nil
	}

	// story rejected
	return false, nil
}

// calculate weighted vote based on user's total category coin balance
func weightedVote(
	ctx sdk.Context, accountKeeper auth.AccountKeeper, votes []interface{}) (
	weightedAmount sdk.Int, err sdk.Error) {

	weightedAmount = sdk.ZeroInt()

	for _, vote := range votes {
		v, ok := vote.(app.Voter)
		if !ok {
			return weightedAmount, ErrInvalidVote(v)
		}

		user := accountKeeper.GetAccount(ctx, v.VoteCreator())
		categoryCoins := user.GetCoins().AmountOf(v.AmountDenom())
		weightedAmount = weightedAmount.Add(categoryCoins)
	}

	return weightedAmount, nil
}
