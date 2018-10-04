package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Implements a queue on top of the key-value store. It keeps a list of ids
// for unexpired backings which are checked for maturity on each block tick.

// unexported key for backing queue
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

// BackQueueLen gets the length of the queue
func (k TruKeeper) BackQueueLen(ctx sdk.Context) int {
	return len(k.getBackingQueue(ctx))
}

// ============================================================================

// getBackingQueue gets the StoryQueue from the context
func (k TruKeeper) getBackingQueue(ctx sdk.Context) ts.BackingQueue {
	store := ctx.KVStore(k.backingKey)
	bq := store.Get(keyBackingQueue)
	if bq == nil {
		return ts.BackingQueue{}
	}
	q := &ts.BackingQueue{}
	k.cdc.MustUnmarshalBinary(bq, q)

	return *q
}

// setBackingQueue sets the BackingQueue to the context
func (k TruKeeper) setBackingQueue(ctx sdk.Context, q ts.BackingQueue) {
	store := ctx.KVStore(k.backingKey)
	bsq := k.cdc.MustMarshalBinary(q)
	store.Set(keyBackingQueue, bsq)
}
