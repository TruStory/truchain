package story

import (
	"fmt"

	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewResponseEndBlock is called at the end of every block tick
func (k Keeper) NewResponseEndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processStoryQueue(ctx, k.storyQueue(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// processStoryQueue checks to see if a story has expired
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

	if story.State == Challenged {
		// add story id to challenged story queue
	}

	// move challenges and votes to associate with story id, not game
	// remove game

	// process challenged story queue
	// -- check block time > voting end time (add field to story)
	// -- if so, tally and distribute
	// -- if not, move to next block

	// get rid of pending game list and game queue

	logger.Info("Processing " + story.String())

	if ctx.BlockHeader().Time.Before(story.ExpireTime) {
		// story hasn't expired yet
		// terminate and wait until the next block
		return nil
	}

	logger.Info(fmt.Sprintf("Handling expired story: %d", story.ID))

	storyQueue.Pop()
	story.State = Expired
	k.UpdateStory(ctx, story)

	// check next story
	return k.processStoryQueue(ctx, storyQueue)
}
