package account

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB(t)

	_, publicKey, address, coins := getFakeAppAccountParams()
	createdAppAccount, _ := keeper.CreateAppAccount(ctx, address, coins, publicKey)

	params, jsonErr := ModuleCodec.MarshalJSON(QueryAppAccountParams{
		Address: address,
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryAppAccount),
		Data: params,
	}

	result, sdkErr := queryAppAccount(ctx, query, keeper)
	assert.Nil(t, sdkErr)

	var returnedAppAccount AppAccount
	jsonErr = keeper.codec.UnmarshalJSON(result, &returnedAppAccount)
	assert.Nil(t, jsonErr)
	assert.Equal(t, returnedAppAccount.PrimaryAddress(), createdAppAccount.PrimaryAddress())
}

func TestQueryPrimaryAccount_Success(t *testing.T) {
	ctx, keeper := mockDB(t)

	_, publicKey, address, coins := getFakeAppAccountParams()
	keeper.CreateAppAccount(ctx, address, coins, publicKey)

	params, jsonErr := ModuleCodec.MarshalJSON(QueryAppAccountParams{
		Address: address,
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryPrimaryAccount),
		Data: params,
	}

	result, sdkErr := queryPrimaryAccount(ctx, query, keeper)
	assert.NoError(t, sdkErr)

	var returnedPrimaryAccount PrimaryAccount
	jsonErr = keeper.codec.UnmarshalJSON(result, &returnedPrimaryAccount)
	assert.NoError(t, jsonErr)
	assert.Equal(t, publicKey, returnedPrimaryAccount.GetPubKey())
}

func TestQueryAppAccounts_Success(t *testing.T) {
	ctx, keeper := mockDB(t)

	_, publicKey, address, coins := getFakeAppAccountParams()
	createdAppAccount, err := keeper.CreateAppAccount(ctx, address, coins, publicKey)
	assert.NoError(t, err)

	_, publicKey2, address2, coins2 := getFakeAppAccountParams()
	_, err = keeper.CreateAppAccount(ctx, address2, coins2, publicKey2)
	assert.NoError(t, err)

	_, publicKey3, address3, coins3 := getFakeAppAccountParams()
	_, err = keeper.CreateAppAccount(ctx, address3, coins3, publicKey3)
	assert.NoError(t, err)

	queryParams := QueryAppAccountsParams{
		Addresses: []sdk.AccAddress{address, address2, address3},
	}
	queryParamsBytes, jsonErr := ModuleCodec.MarshalJSON(queryParams)
	assert.NoError(t, jsonErr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QueryAppAccounts}, "/"),
		Data: queryParamsBytes,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryAppAccounts}, query)
	require.NoError(t, err)

	returnedAppAccounts := make([]AppAccount, 0, len(queryParams.Addresses))

	jsonErr = ModuleCodec.UnmarshalJSON(resBytes, &returnedAppAccounts)
	assert.NoError(t, jsonErr)
	assert.Equal(t, returnedAppAccounts[0].PrimaryAddress(), createdAppAccount.PrimaryAddress())
	assert.Len(t, returnedAppAccounts, 3)
}

func TestQueryAppAccount_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB(t)

	_, _, address, _ := getFakeAppAccountParams()

	params, jsonErr := ModuleCodec.MarshalJSON(QueryAppAccountParams{
		Address: address,
	})
	assert.NoError(t, jsonErr)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryAppAccount),
		Data: params,
	}

	_, sdkErr := queryAppAccount(ctx, query, keeper)
	assert.NotNil(t, sdkErr)
	assert.Equal(t, ErrAppAccountNotFound(address).Code(), sdkErr.Code())
}

func TestQueryParams_Success(t *testing.T) {
	ctx, keeper := mockDB(t)

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
