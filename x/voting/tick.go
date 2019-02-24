package voting

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processVotingStoryQueue(ctx)
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// Recursively process voting story queue to see if voting has ended
func (k Keeper) processVotingStoryQueue(ctx sdk.Context) sdk.Error {
	votingStoryQueue := k.votingStoryQueue(ctx)

	if votingStoryQueue.IsEmpty() {
		return nil
	}

	var storyID int64
	peekErr := votingStoryQueue.Peek(&storyID)
	if peekErr != nil {
		panic(peekErr)
	}

	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return err
	}

	if ctx.BlockHeader().Time.Before(story.VotingEndTime) {
		// no stories to process
		// check again after next block
		return nil
	}

	// only left with voting ended stories
	votingStoryQueue.Pop()

	err = k.verifyStory(ctx, storyID)
	if err != nil {
		return err
	}

	// process next story
	return k.processVotingStoryQueue(ctx)
}
