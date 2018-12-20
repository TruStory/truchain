package backing

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func validEvidence() []string {
	return []string{"http://www.trustory.io"}
}

func TestValidBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := DefaultMsgParams().MinPeriod
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, validStake, validArgument, validCreator, validPeriod, validEvidence())
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "backing", msg.Route())
	assert.Equal(t, "back_story", msg.Type())
	assert.Equal(t, validStake.Amount.String(), msg.Amount.Amount.String())
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInvalidStoryIdBackMsg(t *testing.T) {
	invalidStoryID := int64(-1)
	validStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	validArgument := "valid argument"
	msg := NewBackStoryMsg(invalidStoryID, validStake, validArgument, validCreator, validPeriod, validEvidence())
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(702), err.Code(), err.Error())
}

func TestInvalidAddressBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(100)}
	invalidCreator := sdk.AccAddress([]byte{})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, validStake, validArgument, invalidCreator, validPeriod, validEvidence())
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidStakeBackMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(0)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, invalidStake, validArgument, validCreator, validPeriod, validEvidence())
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(5), err.Code(), err.Error())
}

func TestInValidBackingPeriodBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidPeriod := time.Duration(0 * time.Hour)
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, validStake, validArgument, validCreator, invalidPeriod, validEvidence())
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(901), err.Code(), err.Error())
}

func TestInValidBackingPeriod2BackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidPeriod := time.Duration(366 * 24 * time.Hour)
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, validStake, validArgument, validCreator, invalidPeriod, validEvidence())
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(901), err.Code(), err.Error())
}
