package db

import (
	"testing"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func TestBackingQueue_ErrBackingQueueEmpty(t *testing.T) {
	ctx, _, _, k := mockDB()

	// test empty queue
	_, err := k.BackingQueueHead(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, ts.ErrBackingQueueEmpty().Code(), err.Code(), "should be empty")

	_, err = k.BackingQueuePop(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, ts.ErrBackingQueueEmpty().Code(), err.Code(), "should be empty")

	k.BackingQueuePush(ctx, int64(5))
	_, err = k.BackingQueueHead(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, ts.ErrBackingNotFound(5).Code(), err.Code(), "backing should not be found")
}

func TestBackingQueue(t *testing.T) {
	ctx, ms, _, k := mockDB()

	// create fake backing
	storyID := createFakeStory(ms, k)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := ts.NewBackingParams().MinPeriod
	k.ck.AddCoins(ctx, creator, sdk.Coins{amount})

	// create backings
	backingID, err := k.NewBacking(ctx, storyID, amount, creator, duration)
	_, err = k.GetBacking(ctx, backingID)
	assert.Nil(t, err)

	backingID, err = k.NewBacking(ctx, storyID, amount, creator, duration)
	_, err = k.GetBacking(ctx, backingID)
	assert.Nil(t, err)

	len := k.BackQueueLen(ctx)
	assert.Equal(t, 2, len, "length of queue should be correct")

	backing, _ := k.BackingQueuePop(ctx)
	assert.Equal(t, backing.ID, int64(0), "backing id should match")

	backing, _ = k.BackingQueuePop(ctx)
	assert.Equal(t, backing.ID, int64(1), "backing id should match")

	len = k.BackQueueLen(ctx)
	assert.Equal(t, 0, len, "length of queue should be correct")
}
