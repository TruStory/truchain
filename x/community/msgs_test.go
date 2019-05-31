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
	assert.Equal(t, StoreKey, msg.Route())
	assert.Equal(t, "new_community", msg.Type())
}

func TestMsgNewCommunity_InvalidName(t *testing.T) {
	_, slug, description := getFakeCommunityParams()
	creator := sdk.AccAddress([]byte{1, 2})
	invalidName := "Some really really really long name for a community"

	msg := NewMsgNewCommunity(invalidName, slug, description, creator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestMsgNewCommunity_InvalidSlug(t *testing.T) {
	name, _, description := getFakeCommunityParams()
	creator := sdk.AccAddress([]byte{1, 2})
	invalidSlug := "some-really-really-really-long-name-for-a-community"

	msg := NewMsgNewCommunity(name, invalidSlug, description, creator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestMsgNewCommunity_InvalidDescription(t *testing.T) {
	name, slug, _ := getFakeCommunityParams()
	creator := sdk.AccAddress([]byte{1, 2})
	invalidDescription := "If I could ever think of a really silly day of my life, I would choose the day when I tried fitting in more than 140 chars in a tweet. How silly it was!"

	msg := NewMsgNewCommunity(name, slug, invalidDescription, creator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestMsgNewCommunity_InvalidCreator(t *testing.T) {
	name, slug, description := getFakeCommunityParams()
	invalidCreator := sdk.AccAddress(nil)

	msg := NewMsgNewCommunity(name, slug, description, invalidCreator)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}
