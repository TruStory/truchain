package account

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins := getFakeAppAccountParams()
	createdAppAccount, _ := keeper.CreateAppAccount(ctx, address, coins, publicKey)

	params, jsonErr := json.Marshal(QueryAppAccountParams{
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
	assert.Equal(t, returnedAppAccount.GetAddress(), createdAppAccount.GetAddress())
}

func TestQueryAppAccounts_Success(t *testing.T) {
	ctx, keeper := mockDB()

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
	assert.Equal(t, returnedAppAccounts[0].GetAddress(), createdAppAccount.GetAddress())
	assert.Len(t, returnedAppAccounts, 3)
}

func TestQueryAppAccount_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	_, _, address, _ := getFakeAppAccountParams()

	params, jsonErr := json.Marshal(QueryAppAccountParams{
		Address: address,
	})
	assert.NoError(t, jsonErr)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryAppAccount),
		Data: params,
	}

	_, sdkErr := queryAppAccount(ctx, query, keeper)
	assert.NotNil(t, sdkErr)
	t.Log(sdkErr)
	assert.Equal(t, ErrAppAccountNotFound(address).Code(), sdkErr.Code())
}
