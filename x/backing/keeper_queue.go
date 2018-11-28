package backing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Implements a queue on top of the key-value store. It keeps a list of ids
// for unexpired backings which are checked for maturity on each block tick.

// QueueHead returns the head of the FIFO queue
func (k Keeper) QueueHead(ctx sdk.Context) (backing Backing, err sdk.Error) {
	q := k.getQueue(ctx)
	if q.IsEmpty() {
		return backing, ErrQueueEmpty()
	}
	if backing, err = k.Backing(ctx, q[0]); err != nil {
		return
	}
	return
}

// QueuePop pops the head from the backing queue
func (k Keeper) QueuePop(ctx sdk.Context) (backing Backing, err sdk.Error) {
	q := k.getQueue(ctx)
	if q.IsEmpty() {
		return backing, ErrQueueEmpty()
	}
	headElement, tailQueue := q[0], q[1:]
	k.setQueue(ctx, tailQueue)
	if backing, err = k.Backing(ctx, headElement); err != nil {
		return
	}

	return
}

// QueuePush pushes a backing to the tail of the FIFO queue
func (k Keeper) QueuePush(ctx sdk.Context, id int64) {
	q := k.getQueue(ctx)
	k.setQueue(ctx, append(q, id))
}

// QueueLen gets the length of the queue
func (k Keeper) QueueLen(ctx sdk.Context) int {
	return len(k.getQueue(ctx))
}

// ============================================================================

// geQueue gets the queue from the context
func (k Keeper) getQueue(ctx sdk.Context) (q Queue) {
	store := k.GetStore(ctx)
	bq := store.Get(getQueueKey(k))
	if bq == nil {
		return
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(bq, &q)

	return
}

// setQueue sets the Queue to the context
func (k Keeper) setQueue(ctx sdk.Context, q Queue) {
	store := k.GetStore(ctx)
	bsq := k.GetCodec().MustMarshalBinaryLengthPrefixed(q)
	store.Set(getQueueKey(k), bsq)
}

// getQueueKey returns "backings:unexpired:queue"
func getQueueKey(k Keeper) []byte {
	return []byte(fmt.Sprintf(
		"%s:%s:%s",
		k.GetStoreKey().Name(),
		"unexpired",
		"queue"))
}
