package backing

import (
	"fmt"
	"net/url"
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
	argument := "cool story brew"
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{3, 4})
	duration := DefaultMsgParams().MinPeriod
	testURL, _ := url.Parse("http://www.trustory.io")
	evidence := []url.URL{*testURL}

	bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	// create backings
	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator1, duration, evidence)
	fmt.Println(backingID)
	_, err := bk.Backing(ctx, backingID)
	assert.Nil(t, err)

	backingID, _ = bk.Create(ctx, storyID, amount, argument, creator2, duration, evidence)
	_, err = bk.Backing(ctx, backingID)
	assert.Nil(t, err)

	len := bk.QueueLen(ctx)
	assert.Equal(t, 2, len)

	backing, _ := bk.QueuePop(ctx)
	assert.Equal(t, int64(1), backing.ID())

	backing, _ = bk.QueuePop(ctx)
	assert.Equal(t, int64(2), backing.ID())

	len = bk.QueueLen(ctx)
	assert.Equal(t, 0, len)
}
