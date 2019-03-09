package voting

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processChallengedStoryQueue(ctx)
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// Recursively process voting story queue to see if voting has ended
func (k Keeper) processChallengedStoryQueue(ctx sdk.Context) sdk.Error {
	logger := ctx.Logger().With("module", StoreKey)

	challengedStoryQueue := k.challengedStoryQueue(ctx)

	if challengedStoryQueue.IsEmpty() {
		return nil
	}

	logger.Info("Processing challenged story queue...")

	var storyID int64
	peekErr := challengedStoryQueue.Peek(&storyID)
	if peekErr != nil {
		panic(peekErr)
	}

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	logger.Info("Checking story " + story.String())

	if ctx.BlockHeader().Time.Before(story.VotingEndTime) {
		// no stories to process
		// check again after next block
		logger.Info("Current block time: " + ctx.BlockHeader().Time.String())
		logger.Info(fmt.Sprintf("Story %s is still in voting period", story))

		return nil
	}

	// only left with voting ended stories
	challengedStoryQueue.Pop()

	err = k.verifyStory(ctx, storyID)
	if err != nil {
		return err
	}

	// process next story
	return k.processChallengedStoryQueue(ctx)
}

// tally votes and distribute rewards
func (k Keeper) verifyStory(ctx sdk.Context, storyID int64) sdk.Error {
	logger := ctx.Logger().With("module", StoreKey)

	logger.Info(fmt.Sprintf("Verifying story id: %d...", storyID))

	// tally backings, challenges, and votes
	votes, err := k.tally(ctx, storyID)
	if err != nil {
		return err
	}

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	credDenom, err := k.storyKeeper.CategoryDenom(ctx, storyID)
	if err != nil {
		return err
	}

	// check if story was confirmed
	confirmed, err := k.confirmStory(ctx, votes, credDenom, storyID)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Story confirmed: %t", confirmed))

	// calculate reward pool
	rewardPool, err := k.rewardPool(ctx, votes, confirmed, story.CategoryID)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Reward pool: %v", rewardPool))

	// distribute rewards
	err = k.distributeRewards(ctx, rewardPool, votes, confirmed, story.CategoryID)
	if err != nil {
		return err
	}

	// update story state
	err = k.storyKeeper.EndVotingPeriod(ctx, storyID, confirmed)
	if err != nil {
		return err
	}

	return nil
}
