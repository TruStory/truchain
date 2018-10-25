package backing

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func TestQueue(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()

	// create fake backing
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := NewParams().MinPeriod
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	// create backings
	backingID, _ := bk.NewBacking(ctx, storyID, amount, creator, duration)
	_, err := bk.GetBacking(ctx, backingID)
	assert.Nil(t, err)

	backingID, _ = bk.NewBacking(ctx, storyID, amount, creator, duration)
	_, err = bk.GetBacking(ctx, backingID)
	assert.Nil(t, err)

	len := bk.QueueLen(ctx)
	assert.Equal(t, 2, len, "length of queue should be correct")

	backing, _ := bk.QueuePop(ctx)
	assert.Equal(t, backing.ID, int64(1), "backing id should match")

	backing, _ = bk.QueuePop(ctx)
	assert.Equal(t, backing.ID, int64(2), "backing id should match")

	len = bk.QueueLen(ctx)
	assert.Equal(t, 0, len, "length of queue should be correct")
}
