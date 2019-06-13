package auth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()
	createdAppAccount := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0)

	params, jsonErr := json.Marshal(QueryAppAccountParams{
		Address: address,
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: "/custom/auth/address",
		Data: params,
	}

	result, sdkErr := queryAppAccount(ctx, query, keeper)
	assert.Nil(t, sdkErr)

	var returnedAppAccount AppAccount
	jsonErr = json.Unmarshal(result, &returnedAppAccount)
	assert.Nil(t, jsonErr)
	assert.Equal(t, returnedAppAccount.BaseAccount, createdAppAccount.BaseAccount)
}

func TestQueryCommunity_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	_, _, address, _, _ := getFakeAppAccountParams()

	params, jsonErr := json.Marshal(QueryAppAccountParams{
		Address: address,
	})
	assert.Nil(t, jsonErr)

	query := abci.RequestQuery{
		Path: "/custom/auth/address",
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
	first := keeper.NewAppAccount(ctx, address, coins, publicKey, 0, 0)

	_, publicKey2, address2, coins2, _ := getFakeAppAccountParams()
	another := keeper.NewAppAccount(ctx, address2, coins2, publicKey2, 0, 0)

	result, sdkErr := queryAppAccounts(ctx, keeper)
	assert.Nil(t, sdkErr)

	var appAccounts []AppAccount
	jsonErr := json.Unmarshal(result, &appAccounts)
	t.Log(appAccounts)
	assert.Nil(t, jsonErr)
	assert.Equal(t, appAccounts[0].BaseAccount.Address, first.BaseAccount.Address)
	assert.Equal(t, appAccounts[1].BaseAccount.Address, another.BaseAccount.Address)
}
