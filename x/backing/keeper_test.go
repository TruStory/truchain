package backing

import (
	"math"
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
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
	backingID, _ := bk.Create(ctx, storyID, amount, 0, argument, creator)

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

	bk.Create(ctx, storyID, amount, 0, argument, creator)
	bk.Create(ctx, storyID, amount, 0, argument, creator2)

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

	bk.Create(ctx, storyID, amount, 0, argument, creator)

	backing, _ := bk.BackingByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), backing.ID())
}

func TestTotalBacking(t *testing.T) {
	ctx, k, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, 0, argument, creator)

	creator2 := sdk.AccAddress([]byte{2, 3})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, 0, argument, creator2)

	total, _ := k.TotalBackingAmount(ctx, storyID)

	assert.Equal(t, "10000000trusteak", total.String())
}

func TestNewBacking_ErrInsufficientFunds(t *testing.T) {
	ctx, bk, sk, ck, _, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	_, err := bk.Create(ctx, storyID, amount, 0, argument, creator)
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

	backingID, _ := bk.Create(ctx, storyID, amount, 0, argument, creator)
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

	backingID, _ := bk.Create(ctx, storyID, amount, 0, argument, creator)
	assert.NotNil(t, backingID)
	_, err := bk.Create(ctx, storyID, amount, 0, argument, creator)
	assert.Equal(t, ErrDuplicate(storyID, creator).Code(), err.Code())
}

func Test_BackersByStoryID(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{1, 2, 3, 4})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	_, err := bk.Create(ctx, storyID, amount, 0, argument, creator)
	assert.Nil(t, err)

	_, err = bk.Create(ctx, storyID, amount, 0, argument, creator2)
	assert.Nil(t, err)

	backers, err := bk.BackersByStoryID(ctx, storyID)
	assert.Nil(t, err)
	assert.Subset(t, []sdk.Address{creator, creator2}, backers)
}
