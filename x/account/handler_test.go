package account

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestHandleMsgRegisterKey(t *testing.T) {
	ctx, keeper := mockDB(t)
	handler := NewHandler(keeper)
	assert.NotNil(t, handler) // assert handler is present

	_, publicKey, address, coins := getFakeAppAccountParams()

	registrar := keeper.GetParams(ctx).Registrar

	msg := NewMsgRegisterKey(registrar, address, publicKey, "secp256k1", coins)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	var appAccount AppAccount
	err := keeper.codec.UnmarshalJSON(result.Data, &appAccount)
	assert.NoError(t, err)

	acc, err := keeper.PrimaryAccount(ctx, address)
	assert.NoError(t, err)
	assert.Equal(t, acc.GetPubKey(), publicKey)
}

func TestByzantineMsg(t *testing.T) {
	ctx, keeper := mockDB(t)

	handler := NewHandler(keeper)
	assert.NotNil(t, handler)

	res := handler(ctx, nil)
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)
}
