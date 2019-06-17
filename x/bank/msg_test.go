package bank

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgPayReward(t *testing.T) {
	brokerAddress := sdk.AccAddress([]byte("from"))
	recipientAddr := sdk.AccAddress([]byte("to"))
	reward := sdk.NewInt64Coin("mydenom", 10)
	msg := NewMsgPayReward(brokerAddress, recipientAddr, reward, 1)
	assert.Equal(t, msg.Route(), RouterKey)
	assert.Equal(t, msg.Type(), "pay_reward")
	assert.Equal(t, []sdk.AccAddress{brokerAddress}, msg.GetSigners())

}

func TestMsgPayReward_ValidateBasic(t *testing.T) {
	brokerAddress := sdk.AccAddress([]byte("from"))
	recipientAddr := sdk.AccAddress([]byte("to"))
	invalidAddress := sdk.AccAddress([]byte(""))
	reward := sdk.NewInt64Coin("mydenom", 10)
	msg := NewMsgPayReward(invalidAddress, recipientAddr, reward, 1)
	err := msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidAddress, err.Code())

	msg = NewMsgPayReward(brokerAddress, invalidAddress, reward, 1)
	err = msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidAddress, err.Code())

	msg = NewMsgPayReward(brokerAddress, recipientAddr, sdk.Coin{Denom: "mydenom", Amount: sdk.NewInt(-1)}, 1)
	err = msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidCoins, err.Code())

	msg = NewMsgPayReward(brokerAddress, recipientAddr, sdk.NewInt64Coin("mydenom", 0), 1)
	err = msg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidCoins, err.Code())
}
