package trubank

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidPayRewardMsg(t *testing.T) {
	validCreator := sdk.AccAddress([]byte{1, 2})
	validRecipient := sdk.AccAddress([]byte{1, 2})
	validReward := sdk.Coin{Denom: "trusteak", Amount: sdk.NewInt(1)}
	validID := int64(1)
	msg := NewPayRewardMsg(validCreator, validRecipient, validReward, validID)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "trubank", msg.Route())
	assert.Equal(t, "pay_reward", msg.Type())
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInvalidCreatorPayRewardMsg(t *testing.T) {
	invalidCreator := sdk.AccAddress([]byte{})
	validRecipient := sdk.AccAddress([]byte{1, 2})
	validReward := sdk.Coin{Denom: "trusteak", Amount: sdk.NewInt(1)}
	validID := int64(1)
	msg := NewPayRewardMsg(invalidCreator, validRecipient, validReward, validID)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInvalidRecipientPayRewardMsg(t *testing.T) {
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidRecipient := sdk.AccAddress([]byte{})
	validReward := sdk.Coin{Denom: "trusteak", Amount: sdk.NewInt(1)}
	validID := int64(1)
	msg := NewPayRewardMsg(validCreator, invalidRecipient, validReward, validID)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInvalidAmountPayRewardMsg(t *testing.T) {
	validCreator := sdk.AccAddress([]byte{1, 2})
	validRecipient := sdk.AccAddress([]byte{1, 2})
	invalidReward := sdk.Coin{Denom: "trusteak", Amount: sdk.NewInt(0)}
	validID := int64(1)
	msg := NewPayRewardMsg(validCreator, validRecipient, invalidReward, validID)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(5), err.Code(), err.Error())
}

func TestInvalidDenomPayRewardMsg(t *testing.T) {
	validCreator := sdk.AccAddress([]byte{1, 2})
	validRecipient := sdk.AccAddress([]byte{1, 2})
	invalidReward := sdk.Coin{Denom: "fakecoin", Amount: sdk.NewInt(1)}
	validID := int64(1)
	msg := NewPayRewardMsg(validCreator, validRecipient, invalidReward, validID)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(11), err.Code(), err.Error())
}
