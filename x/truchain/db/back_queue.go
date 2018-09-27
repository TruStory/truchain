// Implements a queue on top of the Cosmos key-value store.
// Useful for managing lists of data, like stories in-progress.

package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var keyBackingQueue = []byte("backings:queue")

// BackingQueueHead returns the head of the FIFO queue
func (k TruKeeper) BackingQueueHead(ctx sdk.Context) (ts.Backing, sdk.Error) {
	q := k.getBackingQueue(ctx)
	if q.IsEmpty() {
		return ts.Backing{}, ts.ErrBackingQueueEmpty()
	}
	backing, err := k.GetBacking(ctx, q[0])
	if err != nil {
		return ts.Backing{}, err
	}
	return backing, nil
}

// BackingQueuePop pops the head from the backing queue
func (k TruKeeper) BackingQueuePop(ctx sdk.Context) (ts.Backing, sdk.Error) {
	q := k.getBackingQueue(ctx)
	if q.IsEmpty() {
		return ts.Backing{}, ts.ErrBackingQueueEmpty()
	}
	headElement, tailQueue := q[0], q[1:]
	k.setBackingQueue(ctx, tailQueue)
	backing, err := k.GetBacking(ctx, headElement)
	if err != nil {
		return ts.Backing{}, err
	}

	return backing, nil
}

// BackingQueuePush pushes a backing to the tail of the FIFO queue
func (k TruKeeper) BackingQueuePush(ctx sdk.Context, id int64) {
	q := k.getBackingQueue(ctx)
	q = append(q, id)
	k.setBackingQueue(ctx, q)
}

// ============================================================================

// getBackingQueue gets the StoryQueue from the context
func (k TruKeeper) getBackingQueue(ctx sdk.Context) ts.BackingQueue {
	store := ctx.KVStore(k.backingKey)
	bq := store.Get(keyBackingQueue)
	// if queue is empty, create new one
	if bq == nil {
		return ts.BackingQueue{}
	}
	// unmarshal bytes to array
	q := &ts.BackingQueue{}
	k.cdc.MustUnmarshalBinary(bq, q)

	return *q
}

// setBackingQueue sets the BackingQueue to the context
func (k TruKeeper) setBackingQueue(ctx sdk.Context, storyQueue ts.BackingQueue) {
	store := ctx.KVStore(k.storyKey)
	bsq := k.cdc.MustMarshalBinary(storyQueue)
	store.Set(keyBackingQueue, bsq)
}
