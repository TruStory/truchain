package backing

import (
	"math"
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var fiver = sdk.Coin{
	Amount: sdk.NewInt(5),
	Denom:  app.StakeDenom,
}

func Test_key(t *testing.T) {
	_, bk, _, _, _, _ := mockDB()

	bz1 := bk.GetIDKey(5)
	bz2 := bk.GetIDKey(math.MaxInt64)

	assert.Equal(t, "backings:id:5", string(bz1), "should generate valid key")
	assert.Equal(t, "backings:id:9223372036854775807", string(bz2), "should generate valid key")
}

func TestGetBacking_ErrBackingNotFound(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()
	id := int64(5)

	_, err := bk.Backing(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound(id).Code(), err.Code(), "Should get error")
}

func TestGetBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew.."
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator, false)

	b, err := bk.Backing(ctx, backingID)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), b.ID())
}

func TestBackingsByStoryID(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	creator2 := sdk.AccAddress([]byte{2, 3})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	bk.Create(ctx, storyID, amount, argument, creator, false)
	bk.Create(ctx, storyID, amount, argument, creator2, false)

	backings, _ := bk.BackingsByStoryID(ctx, storyID)
	assert.Equal(t, 2, len(backings))
}

func TestBackingsByStoryIDAndCreator(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	bk.Create(ctx, storyID, amount, argument, creator, false)

	backing, _ := bk.BackingByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), backing.ID())
}

func TestTally(t *testing.T) {
	ctx, k, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator, false)

	creator2 := sdk.AccAddress([]byte{2, 3})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator2, false)

	yes, _, _ := k.Tally(ctx, storyID)

	assert.Equal(t, 2, len(yes))
}

func TestTotalBacking(t *testing.T) {
	ctx, k, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator, false)

	creator2 := sdk.AccAddress([]byte{2, 3})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator2, false)

	total, _ := k.TotalBackingAmount(ctx, storyID)

	assert.Equal(t, "10000000trusteak", total.String())
}

func TestNewBacking_ErrInsufficientFunds(t *testing.T) {
	ctx, bk, sk, ck, _, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	_, err := bk.Create(ctx, storyID, amount, argument, creator, false)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInsufficientFunds("blah").Code(), err.Code(), "Should get error")
}

func TestNewBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator, false)
	assert.NotNil(t, backingID)
}

func TestDuplicateBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator, false)
	assert.NotNil(t, backingID)
	_, err := bk.Create(ctx, storyID, amount, argument, creator, false)
	assert.Equal(t, ErrDuplicate(storyID, creator).Code(), err.Code())
}

func Test_BackingDelete(t *testing.T) {
	ctx, k, storyKeeper, ck, bankKeeper, am := mockDB()
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now()})
	storyID := createFakeStory(ctx, storyKeeper, ck)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(10000000000))
	argument := "test argument right here"
	totalCoins := sdk.Coins{sdk.NewCoin(app.StakeDenom, sdk.NewInt(2000000000000))}
	backer := createFakeFundedAccount(ctx, am, totalCoins)

	id, err := k.Create(
		ctx, storyID, amount, argument, backer, false)
	assert.NoError(t, err)

	backing, err := k.Backing(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, backing.Vote)
	assert.Equal(t, totalCoins.Minus(sdk.Coins{amount}), bankKeeper.GetCoins(ctx, backer), "coins should have been deducted")

	k.Delete(ctx, backing)

	_, err = k.Backing(ctx, id)
	assert.Equal(t, CodeNotFound, err.Code())
	assert.Equal(t, totalCoins, bankKeeper.GetCoins(ctx, backer), "coins should have been added back")
}
