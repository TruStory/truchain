package slashing

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgSlashArgument_Success(t *testing.T) {
	msg := NewMsgSlashArgument(1, SlashTypeUnhelpful, SlashReasonFocusedOnPerson, "", sdk.AccAddress([]byte{1, 2}))
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgSlashArgument, msg.Type())
}

func TestMsgSlashArgument_InvalidCreator(t *testing.T) {
	msg := NewMsgSlashArgument(1, SlashTypeUnhelpful, SlashReasonFocusedOnPerson, "", sdk.AccAddress(nil))
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}

func TestMsgSlashArgument_InvalidSlashReason(t *testing.T) {
	msg := NewMsgSlashArgument(1, SlashTypeUnhelpful, SlashReasonOther, "", sdk.AccAddress([]byte{1, 2}))
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidSlashReason("").Code(), err.Code())
}
