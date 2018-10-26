package backing

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	tags := bk.NewResponseEndBlock(ctx)
	assert.Nil(t, tags)
}

func Test_processEarnings_QueueEmpty(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	err := processBacking(ctx, bk)
	assert.Nil(t, err)
}

func Test_processEarnings_UnexpiredBackings(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bk.NewBacking(ctx, storyID, amount, creator, duration)
	bk.NewBacking(ctx, storyID, amount, creator, duration)
	bk.NewBacking(ctx, storyID, amount, creator, duration)

	err := processBacking(ctx, bk)
	assert.Nil(t, err)
}

func Test_processEarnings_ExpiredBackings(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 4 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bk.NewBacking(ctx, storyID, amount, creator, duration)
	bk.NewBacking(ctx, storyID, amount, creator, duration)
	bk.NewBacking(ctx, storyID, amount, creator, duration)
	bk.NewBacking(ctx, storyID, amount, creator, duration)
	bk.NewBacking(ctx, storyID, amount, creator, duration)

	// process each backing recursively until queue is empty
	err := processBacking(ctx, bk)
	assert.Nil(t, err)
}

func Test_distributeEarnings(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	principal, _ := sdk.ParseCoin("5trudex")
	interest, _ := sdk.ParseCoin("2trudex")
	expires := time.Now().Add(24 * time.Hour)
	params := DefaultParams()
	duration := 24 * time.Hour
	creator := sdk.AccAddress([]byte{1, 2})

	backing := Backing{
		int64(1),
		int64(5),
		principal,
		interest,
		expires,
		params,
		duration,
		creator,
	}

	err := distributeEarnings(ctx, bk, backing)
	assert.Nil(t, err)
}
