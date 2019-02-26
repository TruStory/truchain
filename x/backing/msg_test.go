package backing

import (
	"testing"
	"time"

	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := 7 * time.Hour
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, validStake, validArgument, validCreator, validPeriod)
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
	msg := NewBackStoryMsg(invalidStoryID, validStake, validArgument, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, story.CodeInvalidStoryID, err.Code(), err.Error())
}

func TestInvalidAddressBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(100)}
	invalidCreator := sdk.AccAddress([]byte{})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, validStake, validArgument, invalidCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeInvalidAddress, err.Code(), err.Error())
}

func TestInValidStakeBackMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidStake := sdk.Coin{Denom: "trustake", Amount: sdk.NewInt(0)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	validArgument := "valid argument"
	msg := NewBackStoryMsg(validStoryID, invalidStake, validArgument, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeInsufficientFunds, err.Code(), err.Error())
}
