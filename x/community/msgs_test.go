package community

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgNewCommunity_Success(t *testing.T) {
	name, id, description := getFakeCommunityParams()
	creator := sdk.AccAddress([]byte{1, 2})

	msg := NewMsgNewCommunity(name, id, description, creator)
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgNewCommunity, msg.Type())
}

func TestMsgNewCommunity_InvalidCreator(t *testing.T) {
	name, id, description := getFakeCommunityParams()
	invalidCreator := sdk.AccAddress(nil)

	msg := NewMsgNewCommunity(name, id, description, invalidCreator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}

func TestMsgUpdateParams_Success(t *testing.T) {
	updates := Params{
		MinIDLength: 20,
	}
	updatedFields := []string{"min_id_length"}
	updater := sdk.AccAddress([]byte{1, 2})
	msg := NewMsgUpdateParams(updates, updatedFields, updater)
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, msg.Updates, updates)
	assert.Equal(t, msg.Updater, updater)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgUpdateParams, msg.Type())
}
