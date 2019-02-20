package game

import (
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

// recursively check for meeting quorum and challenge threshold
func (k Keeper) checkStories(ctx sdk.Context) sdk.Error {
	// logger := ctx.Logger().With("module", "expiration")

	// expiredStoryQueue := k.expiredStoryQueue(ctx)

	// if expiredStoryQueue.IsEmpty() {
	// 	// done processing all expired stories
	// 	// terminate
	// 	return nil
	// }

	// var storyID int64
	// if err := expiredStoryQueue.Peek(&storyID); err != nil {
	// 	panic(err)
	// }
	// logger.Info(fmt.Sprintf("Handling expired story id: %d", storyID))

	// expiredStoryQueue.Pop()

	// err := k.distributeEarningsToBackers(ctx, storyID)
	// if err != nil {
	// 	return err
	// }

	// err = k.returnFundsToChallengers(ctx, storyID)
	// if err != nil {
	// 	return err
	// }

	// handle next expired story
	return k.checkStories(ctx)
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
