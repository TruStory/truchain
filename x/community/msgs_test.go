package community

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgNewCommunity_Success(t *testing.T) {
	name, slug, description := getFakeCommunityParams()
	creator := sdk.AccAddress([]byte{1, 2})

	msg := NewMsgNewCommunity(name, slug, description, creator)
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgNewCommunity, msg.Type())
}

func TestMsgNewCommunity_InvalidCreator(t *testing.T) {
	name, slug, description := getFakeCommunityParams()
	invalidCreator := sdk.AccAddress(nil)

	msg := NewMsgNewCommunity(name, slug, description, invalidCreator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}
