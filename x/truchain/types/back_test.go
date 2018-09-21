package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewBackStoryMsg(validStoryID, validStake, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "BackStory", msg.Type())
	assert.Equal(
		t,
		`{"amount":{"amount":"100","denom":"trusomecoin"},"creator":"cosmosaccaddr1qypq8zs0ka","duration":259200000000000,"story_id":1}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInvalidStoryIdBackMsg(t *testing.T) {
	invalidStoryID := int64(-1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewBackStoryMsg(invalidStoryID, validStake, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(703), err.Code(), err.Error())
}

func TestInvalidAddressBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	invalidCreator := sdk.AccAddress([]byte{})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewBackStoryMsg(validStoryID, validStake, invalidCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidStakeBackMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(0)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewBackStoryMsg(validStoryID, invalidStake, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(705), err.Code(), err.Error())
}

func TestInValidBondPeriodBackMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidPeriod := time.Duration(0 * time.Hour)
	msg := NewBackStoryMsg(validStoryID, validStake, validCreator, invalidPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(706), err.Code(), err.Error())
}
