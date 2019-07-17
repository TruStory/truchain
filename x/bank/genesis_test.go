package bank

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestDefaultGenesisState(t *testing.T) {
	state := DefaultGenesisState()
	assert.Len(t, state.Transactions, 0)
}

func TestInitGenesis(t *testing.T) {
	ctx, keeper, _ := mockDB()
	_, _, rewardAddr := keyPubAddr()
	_, _, appAccountAddr := keyPubAddr()
	params := Params{
		RewardBrokerAddress: rewardAddr,
	}

	regTx := Transaction{
		ID:                1,
		Type:              TransactionGift,
		AppAccountAddress: appAccountAddr,
		ReferenceID:       0,
		Amount:            sdk.NewInt64Coin("mydenom", 300),
		CreatedTime:       ctx.BlockHeader().Time,
	}

	backTx := Transaction{
		ID:                2,
		Type:              TransactionBacking,
		AppAccountAddress: appAccountAddr,
		ReferenceID:       1,
		Amount:            sdk.NewInt64Coin("mydenom", 50),
		CreatedTime:       ctx.BlockHeader().Time,
	}
	transactions := []Transaction{regTx, backTx}
	genesisState := NewGenesisState(params, transactions)
	InitGenesis(ctx, keeper, genesisState)
	actualGenesis := ExportGenesis(ctx, keeper)
	assert.Equal(t, genesisState, actualGenesis)

	err := ValidateGenesis(genesisState)
	assert.NoError(t, err)

	err = ValidateGenesis(GenesisState{})
	assert.Error(t, err)

	// test association list is imported
	accountTxs := keeper.TransactionsByAddress(ctx, appAccountAddr)
	assert.Equal(t, transactions, accountTxs)

}
