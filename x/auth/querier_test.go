package auth

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()
	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey)

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
	assert.Equal(t, returnedAppAccount.BaseAccount, createdAppAccount.BaseAccount)
}

func TestQueryAppAccount_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	_, _, address, _, _ := getFakeAppAccountParams()

	params, jsonErr := json.Marshal(QueryAppAccountParams{
		Address: address,
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryAppAccount),
		Data: params,
	}

	_, sdkErr := queryAppAccount(ctx, query, keeper)
	assert.NotNil(t, sdkErr)
	t.Log(sdkErr)
	assert.Equal(t, ErrAppAccountNotFound(address).Code(), sdkErr.Code())
}

func TestQueryAppAccounts_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()
	_ = keeper.NewAppAccount(ctx, address, coins, publicKey)

	_, publicKey2, address2, coins2, _ := getFakeAppAccountParams()
	_ = keeper.NewAppAccount(ctx, address2, coins2, publicKey2)

	result, sdkErr := queryAppAccounts(ctx, keeper)
	assert.Nil(t, sdkErr)

	var appAccounts []AppAccount
	jsonErr := keeper.codec.UnmarshalJSON(result, &appAccounts)
	assert.Nil(t, jsonErr)
	assert.Len(t, appAccounts, 2)
}
