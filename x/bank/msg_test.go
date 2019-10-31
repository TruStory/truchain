package bank

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgSendGift(t *testing.T) {
	brokerAddress := sdk.AccAddress([]byte("from"))
	recipientAddr := sdk.AccAddress([]byte("to"))
	reward := sdk.NewInt64Coin("mydenom", 10)
	msg := NewMsgSendGift(brokerAddress, recipientAddr, reward)
	assert.Equal(t, msg.Route(), RouterKey)
	assert.Equal(t, msg.Type(), "send_gift")
	assert.Equal(t, []sdk.AccAddress{brokerAddress}, msg.GetSigners())

}

func TestMsgSendGift_ValidateBasic(t *testing.T) {
	brokerAddress := sdk.AccAddress([]byte("from"))
	recipientAddr := sdk.AccAddress([]byte("to"))
	invalidAddress := sdk.AccAddress([]byte(""))
	reward := sdk.NewInt64Coin("mydenom", 10)
	msg := NewMsgSendGift(invalidAddress, recipientAddr, reward)
	err := msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidAddress, err.Code())

	msg = NewMsgSendGift(brokerAddress, invalidAddress, reward)
	err = msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidAddress, err.Code())

	msg = NewMsgSendGift(brokerAddress, recipientAddr, sdk.Coin{Denom: "mydenom", Amount: sdk.NewInt(-1)})
	err = msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidCoins, err.Code())

	msg = NewMsgSendGift(brokerAddress, recipientAddr, sdk.NewInt64Coin("mydenom", 0))
	err = msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidCoins, err.Code())
}
