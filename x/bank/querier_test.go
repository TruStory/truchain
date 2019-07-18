package bank

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryTransactionsByAddress(t *testing.T) {
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

	querier := NewQuerier(keeper)
	queryParams := QueryTransactionsByAddressParams{
		Address: appAccountAddr,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryTransactionsByAddress}, "/"),
		Data: []byte{},
	}
	// Invalid Params
	bz, err := querier(ctx, []string{QueryTransactionsByAddress}, query)
	assert.Error(t, err)
	assert.Equal(t, err.Code(), ErrorCodeInvalidQueryParams)

	// Valid Query
	query.Data = keeper.codec.MustMarshalJSON(&queryParams)
	bz, err = querier(ctx, []string{QueryTransactionsByAddress}, query)
	assert.NoError(t, err)
	assert.NotNil(t, bz)
	var txs []Transaction
	err2 := keeper.codec.UnmarshalJSON(bz, &txs)
	assert.NoError(t, err2)
	assert.Equal(t, transactions, txs)

	// Invalid query path
	bz, err = querier(ctx, []string{"aquerypath"}, query)
	assert.Error(t, err)
	assert.Equal(t, err.Code(), sdk.CodeUnknownRequest)
}

func TestQueryParams_Success(t *testing.T) {
	ctx, keeper, _ := mockDB()

	onChainParams := keeper.GetParams(ctx)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryParams),
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryParams}, query)
	assert.Nil(t, err)

	var returnedParams Params
	sdkErr := ModuleCodec.UnmarshalJSON(resBytes, &returnedParams)
	assert.Nil(t, sdkErr)
	assert.Equal(t, returnedParams, onChainParams)
}
