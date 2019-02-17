package story

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validSource := "http://shanesbrain.net"
	validStoryType := Default
	validArgument := "argument body"
	msg := NewSubmitStoryMsg(validArgument, validBody, validCategoryID, validCreator, validSource, validStoryType)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Route())
	assert.Equal(t, "submit_story", msg.Type())
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

// func TestValidSubmitStoryUnicodeMsg(t *testing.T) {
// 	validBody := "你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好你好好好好好"
// 	assert.Equal(t, DefaultMsgParams().MaxStoryLength, len([]rune(validBody)))

// 	validCategoryID := int64(1)
// 	validCreator := sdk.AccAddress([]byte{1, 2})
// 	validSource := "http://shanesbrain.net"
// 	validStoryType := Default
// 	validArgument := "argument body"
// 	msg := NewSubmitStoryMsg(validArgument, validBody, validCategoryID, validCreator, validSource, validStoryType)
// 	err := msg.ValidateBasic()
// 	assert.Nil(t, err)
// }

func TestInValidBodySubmitStoryMsg(t *testing.T) {
	invalidBody := ""
	validCategoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validSource := "http://shanesbrain.net"
	validStoryType := Default
	validArgument := "argument"

	msg := NewSubmitStoryMsg(validArgument, invalidBody, validCategoryID, validCreator, validSource, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(701), err.Code(), err.Error())
}

func TestInValidCreatorSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategoryID := int64(1)
	invalidCreator := sdk.AccAddress([]byte{})
	validSource := "http://shanesbrain.net"
	validStoryType := Default
	validArgument := "argument"

	msg := NewSubmitStoryMsg(validArgument, validBody, validCategoryID, invalidCreator, validSource, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestValidFlagStoryMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	msg := NewFlagStoryMsg(validStoryID, validCreator)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Route())
	assert.Equal(t, "flag_story", msg.Type())
	assert.Equal(
		t,
		`{"creator":"cosmos1qypq36vzru","story_id":1}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}
