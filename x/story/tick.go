package story

import (
	"fmt"

	list "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock is called at the end of every block tick
func (k Keeper) EndBlock(ctx sdk.Context) sdk.Tags {
	err := k.processPendingStoryList(ctx, k.pendingStoryList(ctx))
	if err != nil {
		panic(err)
	}

	return sdk.EmptyTags()
}

// ============================================================================

// processPendingStoryList checks to see if a story has expired. It checks the state of
// a story, and pushes it's id to the appropriate queue (voting or expired) to
// be handled later in another end blocker.
func (k Keeper) processPendingStoryList(ctx sdk.Context, pendingStoryList list.List) sdk.Error {
	logger := ctx.Logger().With("module", StoreKey)
	pendingStoriesForRemoval := make([]uint64, 0)
	var currentStoryID int64

	pendingStoryList.Iterate(&currentStoryID, func(index uint64) bool {
		story, err := k.Story(ctx, currentStoryID)
		if err != nil {
			panic(err)
		}
		logger.Info(fmt.Sprintf("Processing %s", story.String()))
		if story.Status == Challenged {
			logger.Info(fmt.Sprintf("Voting began for %s", story.String()))
			k.challengedStoryQueue(ctx).Push(currentStoryID)
			logger.Info(fmt.Sprintf(
				"Pushed story id %d to challenge queue, len %d",
				currentStoryID, k.challengedStoryQueue(ctx).List.Len()))
			pendingStoriesForRemoval = append(pendingStoriesForRemoval, index)
			return false
		}

		if ctx.BlockHeader().Time.Before(story.ExpireTime) {
			return false
		}
		logger.Info(fmt.Sprintf("Handling expired: %d", story.ID))
		pendingStoriesForRemoval = append(pendingStoriesForRemoval, index)
		story.Status = Expired
		k.UpdateStory(ctx, story)

		// Push to the expired story queue, which gets handled in
		// the expiration module. At the end of each block, rewards
		// are distributed to backers, and challengers are returned funds.
		k.expiringStoryQueue(ctx).Push(currentStoryID)
		return false
	})

	for _, pendingItemIndex := range pendingStoriesForRemoval {
		pendingStoryList.Delete(pendingItemIndex)
	}
	return nil
}
