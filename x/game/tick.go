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
		quorum, err := k.quorum(ctx, storyID)
		if err != nil {
			panic(err)
		}

		if quorum < k.minQuorum(ctx) {
			// move to next story id
			return false
		}

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

// quorum returns the total count of backings, challenges, votes
func (k Keeper) quorum(ctx sdk.Context, storyID int64) (total int, err sdk.Error) {
	backings, err := k.backingKeeper.BackingsByStoryID(ctx, storyID)
	if err != nil {
		return
	}

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return
	}

	challenges, err := k.challengeKeeper.ChallengesByStoryID(ctx, story.ID)
	if err != nil {
		return
	}

	tokenVotes, err := k.voteKeeper.TokenVotesByStoryID(ctx, story.ID)
	if err != nil {
		return
	}

	total = len(backings) + len(challenges) + len(tokenVotes)

	return total, nil
}
