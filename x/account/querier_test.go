package account

import (
	"encoding/json"
	"fmt"
	"testing"

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
	assert.Equal(t, returnedAppAccount.BaseAccount, createdAppAccount.BaseAccount)
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
