package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidAddCommentMsg(t *testing.T) {
	validStoryID := int64(1)
	validBody := "This is a test comment on a story."
	validCreator := sdk.AccAddress([]byte{1, 2})
	msg := NewAddCommentMsg(validStoryID, validBody, validCreator)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "AddComment", msg.Type())
	assert.Equal(
		t,
		`{"body":"This is a test comment on a story.","creator":"cosmosaccaddr1qypq8zs0ka","story_id":1}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidStoryIDAddCommentMsg(t *testing.T) {
	invalidStoryID := int64(-1)
	validBody := "This is a test comment on a story."
	validCreator := sdk.AccAddress([]byte{1, 2})
	msg := NewAddCommentMsg(invalidStoryID, validBody, validCreator)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(703), err.Code(), err.Error())
}

func TestInValidCreatorAddCommentMsg(t *testing.T) {
	validStoryID := int64(1)
	validBody := "This is a test comment on a story."
	invalidCreator := sdk.AccAddress([]byte{})
	msg := NewAddCommentMsg(validStoryID, validBody, invalidCreator)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidBodyAddCommentMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidBody := ""
	validCreator := sdk.AccAddress([]byte{1, 2})
	msg := NewAddCommentMsg(validStoryID, invalidBody, validCreator)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(702), err.Code(), err.Error())
}
