package bank

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/bank/exported"
)

func TestKeeper_AddCoin(t *testing.T) {
	ctx, k, _ := mockDB()
	addr := []byte("cosmos123456789")
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*20))
	coins, err := k.AddCoin(ctx,
		sdk.AccAddress(addr),
		amount,
		100,
		TransactionBackingReturned,
	)
	assert.NoError(t, err)

	assert.Equal(t, sdk.NewInt(app.Shanev*20), coins.AmountOf(app.StakeDenom))
	k.IterateUserTransactions(ctx, addr, false, func(tx Transaction) bool {
		assert.Equal(t, uint64(1), tx.ID)
		assert.Equal(t, TransactionBackingReturned, tx.Type)
		assert.Equal(t, amount, tx.Amount)
		assert.Equal(t, uint64(100), tx.ReferenceID)
		assert.Equal(t, "", tx.CommunityID)
		return true
	})
	coins, err = k.AddCoin(ctx,
		sdk.AccAddress(addr),
		sdk.Coin{Denom: app.StakeDenom, Amount: sdk.NewInt(app.Shanev * -1)},
		100,
		TransactionBackingReturned,
	)
	assert.Error(t, err)

}
func TestKeeper_AddSubtractCoinWithCommunity(t *testing.T) {
	ctx, k, _ := mockDB()
	addr := []byte("cosmos123456789")
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*20))
	coins, err := k.AddCoin(ctx,
		sdk.AccAddress(addr),
		amount,
		100,
		TransactionBackingReturned,
		exported.WithCommunityID("cryptotest"),
	)
	assert.NoError(t, err)
	assert.Equal(t, sdk.NewInt(app.Shanev*20), coins.AmountOf(app.StakeDenom))

	_, err = k.SubtractCoin(ctx,
		sdk.AccAddress(addr),
		amount,
		100,
		TransactionBacking,
		exported.WithCommunityID("cryptotestsubtract"),
	)
	assert.NoError(t, err)
	txs := make([]Transaction, 0)
	k.IterateUserTransactions(ctx, addr, false, func(tx Transaction) bool {
		txs = append(txs, tx)
		return false
	})
	assert.Len(t, txs, 2)
	assert.Equal(t, uint64(1), txs[0].ID)
	assert.Equal(t, TransactionBackingReturned, txs[0].Type)
	assert.Equal(t, amount, txs[0].Amount)
	assert.Equal(t, uint64(100), txs[0].ReferenceID)
	assert.Equal(t, "cryptotest", txs[0].CommunityID)

	assert.Equal(t, uint64(2), txs[1].ID)
	assert.Equal(t, TransactionBacking, txs[1].Type)
	assert.Equal(t, amount, txs[1].Amount)
	assert.Equal(t, uint64(100), txs[1].ReferenceID)
	assert.Equal(t, "cryptotestsubtract", txs[1].CommunityID)
}

func TestKeeper_SubtractCoin(t *testing.T) {
	ctx, k, auth := mockDB()

	addr := createFakeFundedAccount(ctx, auth, sdk.NewCoins(sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*30))))

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*10))
	coins, err := k.SubtractCoin(ctx,
		addr,
		amount,
		200,
		TransactionBacking,
	)
	assert.NoError(t, err)
	assert.Equal(t, sdk.NewInt(app.Shanev*20), coins.AmountOf(app.StakeDenom))
	k.IterateUserTransactions(ctx, addr, false, func(tx Transaction) bool {
		assert.Equal(t, uint64(1), tx.ID)
		assert.Equal(t, TransactionBacking, tx.Type)
		assert.Equal(t, amount, tx.Amount)
		assert.Equal(t, uint64(200), tx.ReferenceID)
		assert.Equal(t, "", tx.CommunityID)
		return true
	})

	coins, err = k.SubtractCoin(ctx,
		sdk.AccAddress(addr),
		sdk.Coin{Denom: app.StakeDenom, Amount: sdk.NewInt(app.Shanev * -1)},
		100,
		TransactionBacking,
	)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidCoins, err.Code())

	coins, err = k.SubtractCoin(ctx,
		sdk.AccAddress(addr),
		sdk.Coin{Denom: app.StakeDenom, Amount: sdk.NewInt(app.Shanev * 1)},
		100,
		TransactionBackingReturned,
	)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeInvalidTransactionType, err.Code())
}

func TestKeeper_SubtractMoreThanBalanceCoin(t *testing.T) {
	ctx, k, auth := mockDB()

	balance := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*30))
	addr := createFakeFundedAccount(ctx, auth, sdk.NewCoins(balance))

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*200))
	coins, err := k.SafeSubtractCoin(ctx,
		addr,
		amount,
		200,
		TransactionBacking,
	)
	assert.NoError(t, err)
	assert.Equal(t, "0", coins.AmountOf(app.StakeDenom).String())

	k.IterateUserTransactions(ctx, addr, false, func(tx Transaction) bool {
		assert.Equal(t, uint64(1), tx.ID)
		assert.Equal(t, TransactionBacking, tx.Type)
		assert.Equal(t, balance.String(), tx.Amount.String())
		assert.Equal(t, uint64(200), tx.ReferenceID)
		return true
	})
}

func TestKeeper_TransactionsByAddress(t *testing.T) {
	ctx, k, auth := mockDB()

	addr := createFakeFundedAccount(ctx, auth, sdk.NewCoins(sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*100))))

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*10))

	_, err := k.AddCoin(ctx, addr, amount, 200, TransactionRegistration)
	assert.NoError(t, err)
	_, err = k.SubtractCoin(ctx, addr, amount, 200, TransactionBacking)
	assert.NoError(t, err)
	_, err = k.AddCoin(ctx, addr, amount, 200, TransactionBackingReturned)
	assert.NoError(t, err)
	_, err = k.SubtractCoin(ctx, addr, amount, 200, TransactionUpvote)
	assert.NoError(t, err)
	_, err = k.AddCoin(ctx, addr, amount, 200, TransactionUpvoteReturned)
	assert.NoError(t, err)

	txs := k.TransactionsByAddress(ctx, addr)
	assert.Len(t, txs, 5)
	txTypes := make([]TransactionType, 0)
	for _, tx := range txs {
		txTypes = append(txTypes, tx.Type)
	}
	assert.Equal(t,
		[]TransactionType{TransactionRegistration,
			TransactionBacking, TransactionBackingReturned,
			TransactionUpvote, TransactionUpvoteReturned},
		txTypes)

	txs = k.TransactionsByAddress(ctx, addr, FilterByTransactionType(TransactionRegistration, TransactionUpvote))
	txTypes = make([]TransactionType, 0)
	for _, tx := range txs {
		txTypes = append(txTypes, tx.Type)
	}
	assert.Equal(t,
		[]TransactionType{TransactionRegistration, TransactionUpvote},
		txTypes)

	// Test Reverse
	txs = k.TransactionsByAddress(ctx, addr, SortOrder(SortDesc))
	txTypes = make([]TransactionType, 0)

	assert.Len(t, txs, 5)
	txTypes = make([]TransactionType, 0)
	for _, tx := range txs {
		txTypes = append(txTypes, tx.Type)
	}
	assert.Equal(t,
		[]TransactionType{TransactionUpvoteReturned, TransactionUpvote,
			TransactionBackingReturned, TransactionBacking,
			TransactionRegistration,
		},
		txTypes)

	// Limit && Skip
	txs = k.TransactionsByAddress(ctx, addr, Offset(2), Limit(2))
	assert.Len(t, txs, 2)
	txTypes = make([]TransactionType, 0)
	for _, tx := range txs {
		txTypes = append(txTypes, tx.Type)
	}
	assert.Equal(t,
		[]TransactionType{TransactionBackingReturned, TransactionUpvote},
		txTypes)

	txs = k.TransactionsByAddress(ctx, addr, SortOrder(SortDesc), Offset(3), Limit(2))
	assert.Len(t, txs, 2)
	txTypes = make([]TransactionType, 0)
	for _, tx := range txs {
		txTypes = append(txTypes, tx.Type)
	}
	assert.Equal(t,
		[]TransactionType{TransactionBacking, TransactionRegistration},
		txTypes)

}
