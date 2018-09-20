// Implements a queue on top of the Cosmos key-value store.
// Useful for managing lists of data, like stories in-progress.

package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var keyActiveStoryQueue = []byte("activeStoriesQueue")

// ActiveStoryQueueHead returns the head of the FIFO queue
func (k TruKeeper) ActiveStoryQueueHead(ctx sdk.Context) (ts.Story, sdk.Error) {
	storyQueue := k.getActiveStoryQueue(ctx)
	if storyQueue.IsEmpty() {
		return ts.Story{}, ts.ErrActiveStoryQueueEmpty()
	}
	story, err := k.GetStory(ctx, storyQueue[0])
	if err != nil {
		return ts.Story{}, err
	}
	return story, nil
}

// ActiveStoryQueuePop pops the head from the story queue
func (k TruKeeper) ActiveStoryQueuePop(ctx sdk.Context) (ts.Story, sdk.Error) {
	storyQueue := k.getActiveStoryQueue(ctx)
	if storyQueue.IsEmpty() {
		return ts.Story{}, ts.ErrActiveStoryQueueEmpty()
	}
	headElement, tailStoryQueue := storyQueue[0], storyQueue[1:]
	k.setActiveStoryQueue(ctx, tailStoryQueue)
	story, err := k.GetStory(ctx, headElement)
	if err != nil {
		return ts.Story{}, err
	}

	return story, nil
}

// ActiveStoryQueuePush pushes a story to the tail of the FIFO queue
func (k TruKeeper) ActiveStoryQueuePush(ctx sdk.Context, storyID int64) {
	storyQueue := k.getActiveStoryQueue(ctx)
	storyQueue = append(storyQueue, storyID)
	k.setActiveStoryQueue(ctx, storyQueue)
}

// ============================================================================

// getActiveStoryQueue gets the StoryQueue from the context
func (k TruKeeper) getActiveStoryQueue(ctx sdk.Context) ts.ActiveStoryQueue {
	store := ctx.KVStore(k.storyKey)
	bsq := store.Get(keyActiveStoryQueue)
	// if queue is empty, create new one
	if bsq == nil {
		return ts.ActiveStoryQueue{}
	}
	storyQueue := &ts.ActiveStoryQueue{}
	k.cdc.MustUnmarshalBinary(bsq, storyQueue)

	return *storyQueue
}

// setActiveStoryQueue sets the ActiveStoryQueue to the context
func (k TruKeeper) setActiveStoryQueue(ctx sdk.Context, storyQueue ts.ActiveStoryQueue) {
	store := ctx.KVStore(k.storyKey)
	bsq := k.cdc.MustMarshalBinary(storyQueue)
	store.Set(keyActiveStoryQueue, bsq)
}
