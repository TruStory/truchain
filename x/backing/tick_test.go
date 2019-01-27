package backing

import (
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
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
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)

	err := processBacking(ctx, bk)
	assert.Nil(t, err)
}

func Test_processEarnings_ExpiredBackings(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 4 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)

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

	// create new vote type
	vote := app.Vote{
		ID:        int64(1),
		Amount:    principal,
		Argument:  "",
		Creator:   creator,
		Vote:      true,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	// create new backing type with embedded vote
	backing := Backing{
		vote,
		int64(5),
		interest,
		expires,
		params,
		duration,
	}

	err := distributeEarnings(ctx, bk, backing)
	assert.Nil(t, err)
}
