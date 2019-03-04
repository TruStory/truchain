package story

import (
	"fmt"

	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processPendingStoryQueue(ctx, k.pendingStoryQueue(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// processPendingStoryQueue checks to see if a story has expired. It checks the state of
// a story, and pushes it's id to the appropriate queue (voting or expired) to
// be handled later in another end blocker.
func (k Keeper) processPendingStoryQueue(ctx sdk.Context, pendingStoryQueue queue.Queue) sdk.Error {
	logger := ctx.Logger().With("module", "story")

	if pendingStoryQueue.IsEmpty() {
		return nil
	}

	var storyID int64
	if err := pendingStoryQueue.Peek(&storyID); err != nil {
		panic(err)
	}

	story, err := k.Story(ctx, storyID)
	if err != nil {
		return err
	}

	logger.Info("Processing " + story.String())

	// if the status of the story has changed to challenged,
	// add it to the voting story queue to be handled later
	if story.Status == Challenged {
		logger.Info("Voting begun for " + story.String())
		k.challengedStoryQueue(ctx).Push(storyID)

		// pop and process next story
		pendingStoryQueue.Pop()
		return k.processPendingStoryQueue(ctx, pendingStoryQueue)
	}

	if ctx.BlockHeader().Time.Before(story.ExpireTime) {
		// story hasn't expired yet
		// terminate and wait until the next block
		return nil
	}

	logger.Info(fmt.Sprintf("Handling expired: %d", story.ID))

	pendingStoryQueue.Pop()
	story.Status = Expired
	k.UpdateStory(ctx, story)

	// Push to the expired story queue, which gets handled in
	// the expiration module. At the end of each block, rewards
	// are distributed to backers, and challengers are returned funds.
	k.expiredStoryQueue(ctx).Push(storyID)

	// check next story
	return k.processPendingStoryQueue(ctx, pendingStoryQueue)
}
