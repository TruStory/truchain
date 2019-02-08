package vote

import (
	params "github.com/TruStory/truchain/parameters"
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
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

	credDenom, err := k.storyKeeper.CategoryDenom(ctx, game.StoryID)
	if err != nil {
		return err
	}

	// check if story was confirmed
	confirmed, err := confirmStory(ctx, k.accountKeeper, votes, credDenom)
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
		ctx, k.backingKeeper, k.bankKeeper, rewardPool, votes, confirmed, credDenom)
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
	pool = sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())

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
	ctx sdk.Context,
	backingKeeper backing.WriteKeeper,
	bankKeeper bank.Keeper,
	rewardPool sdk.Coin,
	votes poll,
	confirmed bool,
	denom string) (err sdk.Error) {

	logger := ctx.Logger().With("module", "vote")

	if confirmed {
		err = distributeRewardsConfirmed(
			ctx, backingKeeper, bankKeeper, votes, rewardPool, denom)
	} else {
		err = distributeRewardsRejected(
			ctx, backingKeeper, bankKeeper, votes, rewardPool, denom)
	}
	if err != nil {
		return
	}

	logger.Info("Distributed reward pool of " + rewardPool.String())

	return
}

// determine if a story is confirmed or rejected
func confirmStory(
	ctx sdk.Context, accountKeeper auth.AccountKeeper, votes poll, denom string) (
	confirmed bool, err sdk.Error) {

	// calculate weighted true votes
	trueWeight, err := weightedVote(ctx, accountKeeper, votes.trueVotes, denom)
	if err != nil {
		return confirmed, err
	}

	// calculate weighted false votes
	falseWeight, err := weightedVote(ctx, accountKeeper, votes.falseVotes, denom)
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

// calculate weighted vote based on user's cred balance
func weightedVote(
	ctx sdk.Context, accountKeeper auth.AccountKeeper, votes []app.Voter, denom string) (
	weightedAmount sdk.Int, err sdk.Error) {

	weightedAmount = sdk.ZeroInt()

	// iterate through BCVs
	for _, vote := range votes {
		v, ok := vote.(app.Voter)
		if !ok {
			return weightedAmount, ErrInvalidVote(v)
		}

		user := accountKeeper.GetAccount(ctx, v.Creator())
		coins := user.GetCoins()
		credBalance := coins.AmountOf(denom)
		if credBalance.IsZero() {
			// fix cold-start problem by adding 1 preethi
			// when there is a 0 cred balance so the vote
			// is counted
			credBalance = credBalance.Add(sdk.NewInt(1))
		}

		weightedAmount = weightedAmount.Add(credBalance)
	}

	return weightedAmount, nil
}

// Make sure reward pool is empty (<= 1 coin)
// Makes up for rounding error during division
func checkForEmptyPool(pool sdk.Coin) sdk.Error {
	oneCoin := sdk.NewCoin(params.StakeDenom, sdk.OneInt())
	if !(pool.IsLT(oneCoin) || pool.IsEqual(oneCoin)) {
		return ErrNonEmptyRewardPool(pool)
	}

	return nil
}
