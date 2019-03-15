package trubank

import (
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestAddCoins(t *testing.T) {
	ctx, k, ck := mockDB()
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})

	k.MintAndAddCoin(ctx, creator, cat.ID, 0, 5, sdk.NewInt(1000))
	k.MintAndAddCoin(ctx, creator, cat.ID, 0, 5, sdk.NewInt(1000))
	k.MintAndAddCoin(ctx, creator, cat.ID, 0, 5, sdk.NewInt(1000))

	cat2, _ := ck.GetCategory(ctx, cat.ID)

	assert.Equal(t, "3000trudex", cat2.TotalCred.String())
}

func TestSubtractCoins(t *testing.T) {
	ctx, k, _ := mockDB()
	creator := sdk.AccAddress([]byte{1, 2})

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	_, err := k.AddCoin(ctx, creator, amount, 0, Backing, 0)
	assert.Nil(t, err)

	_, err = k.SubtractCoin(ctx, creator, amount, 0, BackingReturned, 0)
	assert.Nil(t, err)
}

func TestTransactionsByCreator(t *testing.T) {
	ctx, k, _ := mockDB()
	creator := sdk.AccAddress([]byte{1, 2})

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	_, err := k.AddCoin(ctx, creator, amount, 0, Backing, 0)
	assert.Nil(t, err)

	_, err = k.SubtractCoin(ctx, creator, amount, 0, BackingReturned, 0)
	assert.Nil(t, err)

	transactions, err := k.TransactionsByCreator(ctx, creator)
	assert.Nil(t, err)

	assert.NotEmpty(t, transactions)
	assert.Len(t, transactions, 2)
}

func TestTransactionNotAddedIfZeroCoins(t *testing.T) {
	ctx, k, _ := mockDB()
	creator := sdk.AccAddress([]byte{1, 2})

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(0))
	_, err := k.AddCoin(ctx, creator, amount, 0, Backing, 0)
	assert.Nil(t, err)

	_, err = k.SubtractCoin(ctx, creator, amount, 0, BackingReturned, 0)
	assert.Nil(t, err)

	transactions, err := k.TransactionsByCreator(ctx, creator)
	assert.Nil(t, err)

	assert.Len(t, transactions, 0)
}

func TestTransactionNotFound(t *testing.T) {
	ctx, k, _ := mockDB()
	creator := sdk.AccAddress([]byte{1, 2})

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(0))
	_, err := k.AddCoin(ctx, creator, amount, 0, Backing, 0)
	assert.Nil(t, err)

	_, err = k.SubtractCoin(ctx, creator, amount, 0, BackingReturned, 0)
	assert.Nil(t, err)

	_, err = k.Transaction(ctx, 2)

	assert.NotNil(t, err)
	assert.Equal(t, ErrTransactionNotFound(2).Code(), err.Code(), "Should get error")

}
