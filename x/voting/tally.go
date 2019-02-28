package voting

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// tally backings, challenges, and token votes into two true and false slices
func (k Keeper) tally(ctx sdk.Context, storyID int64) (votes poll, err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)
	logger.Info("Tallying votes ...")

	// tally backings
	trueBackings, falseBackings, err := k.backingKeeper.Tally(ctx, storyID)
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
	falseChallenges, err := k.challengeKeeper.Tally(ctx, storyID)
	if err != nil {
		return
	}
	for _, v := range falseChallenges {
		votes.falseVotes = append(votes.falseVotes, v)
	}

	// tally token votes
	trueTokenVotes, falseTokenVotes, err := k.voteKeeper.Tally(ctx, storyID)
	if err != nil {
		return
	}
	for _, v := range trueTokenVotes {
		votes.trueVotes = append(votes.trueVotes, v)
	}
	for _, v := range falseTokenVotes {
		votes.falseVotes = append(votes.falseVotes, v)
	}

	logger.Info(votes.String())

	return votes, nil
}

// determine if a story is confirmed or rejected
func (k Keeper) confirmStory(
	ctx sdk.Context, votes poll, denom string) (confirmed bool, err sdk.Error) {

	// calculate weighted true votes
	trueWeight, err := k.weightedVote(ctx, votes.trueVotes, denom)
	if err != nil {
		return confirmed, err
	}

	// calculate weighted false votes
	falseWeight, err := k.weightedVote(ctx, votes.falseVotes, denom)
	if err != nil {
		return confirmed, err
	}

	// calculate what percent of the total weight is true votes
	totalWeight := trueWeight.Add(falseWeight)
	trueWeightDec := sdk.NewDecFromInt(trueWeight)
	truePercentOfTotal := trueWeightDec.QuoInt(totalWeight)

	// majority weight wins
	if truePercentOfTotal.GTE(k.majorityPercent(ctx)) {
		// story confirmed
		return true, nil
	}

	// story rejected
	return false, nil
}

// calculate weighted vote based on user's cred balance
func (k Keeper) weightedVote(
	ctx sdk.Context, votes []app.Voter, denom string) (
	weightedAmount sdk.Int, err sdk.Error) {

	weightedAmount = sdk.ZeroInt()

	// iterate through BCVs
	for _, vote := range votes {
		v, ok := vote.(app.Voter)
		if !ok {
			return weightedAmount, ErrInvalidVote(v)
		}

		user := k.accountKeeper.GetAccount(ctx, v.Creator())
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
