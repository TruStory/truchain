package backing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_ErrQueueEmpty(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	// test empty queue
	_, err := bk.QueueHead(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, ErrQueueEmpty().Code(), err.Code(), "should be empty")

	_, err = bk.QueuePop(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, ErrQueueEmpty().Code(), err.Code(), "should be empty")

	bk.QueuePush(ctx, int64(5))
	_, err = bk.QueueHead(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound(5).Code(), err.Code(), "backing should not be found")
}

// FIX ME SHANNNNNNNE

// func TestQueue(t *testing.T) {
// 	ctx, bk, sk, ck, bankKeeper, _ := mockDB()

// 	// create fake backing
// 	storyID := createFakeStory(ctx, sk, ck)
// 	amount, _ := sdk.ParseCoin("5trudex")
// 	creator := sdk.AccAddress([]byte{1, 2})
// 	duration := DefaultMsgParams().MinPeriod
// 	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

// 	// create backings
// 	backingID, _ := bk.Create(ctx, storyID, amount, creator, duration)
// 	_, err := bk.Backing(ctx, backingID)
// 	assert.Nil(t, err)

// 	backingID, _ = bk.Create(ctx, storyID, amount, creator, duration)
// 	_, err = bk.Backing(ctx, backingID)
// 	spew.Dump(err)
// 	assert.Nil(t, err)

// 	len := bk.QueueLen(ctx)
// 	assert.Equal(t, 2, len)

// 	backing, _ := bk.QueuePop(ctx)
// 	assert.Equal(t, backing.ID, int64(1))

// 	backing, _ = bk.QueuePop(ctx)
// 	assert.Equal(t, backing.ID, int64(2))

// 	len = bk.QueueLen(ctx)
// 	assert.Equal(t, 0, len)
// }
