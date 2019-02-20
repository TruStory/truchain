package story

import (
	"fmt"

	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processStoryQueue(ctx, k.storyQueue(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// processStoryQueue checks to see if a story has expired. It checks the state of
// a story, and pushes it's id to the appropriate queue (voting or expired) to
// be handled later in another end blocker.
func (k Keeper) processStoryQueue(ctx sdk.Context, storyQueue queue.Queue) sdk.Error {
	logger := ctx.Logger().With("module", "story")

	if storyQueue.IsEmpty() {
		return nil
	}

	var storyID int64
	if err := storyQueue.Peek(&storyID); err != nil {
		panic(err)
	}

	story, err := k.Story(ctx, storyID)
	if err != nil {
		return err
	}

	logger.Info("Processing " + story.String())

	// if the state of the story has changed to voting,
	// add it to the voting story queue to be handled later
	if story.State == Voting {
		logger.Info("Voting begun for " + story.String())
		k.votingStoryQueue(ctx).Push(storyID)

		// pop and process next story
		storyQueue.Pop()
		return k.processStoryQueue(ctx, storyQueue)
	}

	if ctx.BlockHeader().Time.Before(story.ExpireTime) {
		// story hasn't expired yet
		// terminate and wait until the next block
		return nil
	}

	logger.Info(fmt.Sprintf("Handling expired: %d", story.ID))

	storyQueue.Pop()
	story.State = Expired
	k.UpdateStory(ctx, story)

	// Push to the expired story queue, which gets handled in
	// the distribution module. At the end of each block, rewards
	// are distributed to backers, and challengers are returned funds.
	k.expiredStoryQueue(ctx).Push(storyID)

	// check next story
	return k.processStoryQueue(ctx, storyQueue)
}
