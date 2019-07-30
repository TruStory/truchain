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

func TestMsgAddAdmin_Success(t *testing.T) {
	admin := sdk.AccAddress([]byte{1, 2})
	creator := sdk.AccAddress([]byte{3, 4})

	msg := NewMsgAddAdmin(admin, creator)
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgAddAdmin, msg.Type())
}

func TestMsgAddAdmin_InvalidCreator(t *testing.T) {
	admin := sdk.AccAddress([]byte{1, 2})
	creator := sdk.AccAddress(nil)

	msg := NewMsgAddAdmin(admin, creator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}

func TestMsgAddAdmin_InvalidAdmin(t *testing.T) {
	admin := sdk.AccAddress(nil)
	creator := sdk.AccAddress([]byte{3, 4})

	msg := NewMsgAddAdmin(admin, creator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}

func TestMsgRemoveAdmin_Success(t *testing.T) {
	remover := sdk.AccAddress([]byte{3, 4})

	msg := NewMsgRemoveAdmin(remover, remover) // self removing
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgRemoveAdmin, msg.Type())
}

func TestMsgRemoveAdmin_InvalidRemover(t *testing.T) {
	admin := sdk.AccAddress([]byte{1, 2})
	remover := sdk.AccAddress(nil)

	msg := NewMsgRemoveAdmin(admin, remover)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}

func TestMsgRemoveAdmin_InvalidAdmin(t *testing.T) {
	admin := sdk.AccAddress(nil)
	remover := sdk.AccAddress([]byte{3, 4})

	msg := NewMsgRemoveAdmin(admin, remover)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}
