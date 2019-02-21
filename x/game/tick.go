package game

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.checkStories(ctx)
	if err != nil {
		panic(err)
	}
	return sdk.EmptyTags()
}

// Iteratively check for meeting quorum and challenge threshold.
// If found, update state of story to voting. Then it'll get picked up
// by the story end blocker, where it'll be added to the voting story queue.
//
// NOTE: This is bit of an expensive operation, and might not be the
// best thing to do at after each block. Watch this space and optimize
// in the future if needed, like moving this into challenge create.
func (k Keeper) checkStories(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", "game")

	var storyID int64
	k.storyQueue(ctx).List.Iterate(&storyID, func(index uint64) bool {
		backingPool, err := k.backingKeeper.TotalBackingAmount(ctx, storyID)
		if err != nil {
			panic(err)
		}

		challengePool, err := k.challengeKeeper.TotalChallengeAmount(ctx, storyID)
		if err != nil {
			panic(err)
		}

		challengeThreshold := k.challengeThreshold(ctx, backingPool)

		logger.Info(fmt.Sprintf(
			"Backing pool: %s, challenge pool: %s, threshold: %s",
			backingPool, challengePool, challengeThreshold))

		if challengePool.IsGTE(challengeThreshold) {
			err := k.storyKeeper.StartVotingPeriod(ctx, storyID)
			if err != nil {
				panic(err)
			}

			logger.Info(fmt.Sprintf(
				"Challenge threshold and quorum met. Voting started for story %d",
				storyID))
		}

		return false
	})

	return nil
}

func (k Keeper) challengeThreshold(ctx sdk.Context, totalBackingAmount sdk.Coin) sdk.Coin {
	// calculate challenge threshold amount (based on total backings)
	totalBackingDec := sdk.NewDecFromInt(totalBackingAmount.Amount)
	challengeThresholdAmount := totalBackingDec.
		Mul(k.challengeToBackingRatio(ctx)).
		TruncateInt()

	// challenge threshold can't be less than min challenge stake
	minChallengeThreshold := k.minChallengeThreshold(ctx)
	if challengeThresholdAmount.LT(minChallengeThreshold) {
		return sdk.NewCoin(totalBackingAmount.Denom, minChallengeThreshold)
	}

	return sdk.NewCoin(totalBackingAmount.Denom, challengeThresholdAmount)
}
