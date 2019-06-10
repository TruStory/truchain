package auth

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgRegisterKey_Success(t *testing.T) {
	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	msg := NewMsgRegisterKey(address, publicKey, "scep", coins)
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgRegisterKey, msg.Type())
}

func TestMsgNewCommunity_InvalidAddress(t *testing.T) {
	_, publicKey, _, coins, _ := getFakeAppAccountParams()
	invalidAddress := sdk.AccAddress(nil)

	msg := NewMsgRegisterKey(invalidAddress, publicKey, "scep", coins)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}
