package db

import (
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, _, _, k := MockDB()

	res := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, res)
}

func Test_processEarnings_BackingQueueEmpty(t *testing.T) {
	ctx, _, _, k := MockDB()

	err := processBacking(ctx, k)
	assert.Nil(t, err)
}

func Test_processEarnings_UnexpiredBackings(t *testing.T) {
	ctx, ms, _, k := MockDB()
	storyID := CreateFakeStory(ms, k)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	k.ck.AddCoins(ctx, creator, sdk.Coins{amount})
	k.NewBacking(ctx, storyID, amount, creator, duration)
	k.NewBacking(ctx, storyID, amount, creator, duration)
	k.NewBacking(ctx, storyID, amount, creator, duration)

	err := processBacking(ctx, k)
	assert.Nil(t, err)
}

func Test_processEarnings_ExpiredBackings(t *testing.T) {
	ctx, ms, _, k := MockDB()
	storyID := CreateFakeStory(ms, k)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 4 * time.Hour
	k.ck.AddCoins(ctx, creator, sdk.Coins{amount})
	k.NewBacking(ctx, storyID, amount, creator, duration)
	k.NewBacking(ctx, storyID, amount, creator, duration)
	k.NewBacking(ctx, storyID, amount, creator, duration)
	k.NewBacking(ctx, storyID, amount, creator, duration)
	k.NewBacking(ctx, storyID, amount, creator, duration)

	// process each backing recursively until queue is empty
	err := processBacking(ctx, k)
	assert.Nil(t, err)
}

func Test_distributeEarnings(t *testing.T) {
	ctx, _, _, k := MockDB()

	principal, _ := sdk.ParseCoin("5trudex")
	interest, _ := sdk.ParseCoin("2trudex")
	expires := time.Now().Add(24 * time.Hour)
	params := ts.NewBackingParams()
	duration := 24 * time.Hour
	creator := sdk.AccAddress([]byte{1, 2})

	backing := ts.NewBacking(
		int64(1),
		int64(5),
		principal,
		interest,
		expires,
		params,
		duration,
		creator)

	err := distributeEarnings(ctx, k, backing)
	assert.Nil(t, err)
}
