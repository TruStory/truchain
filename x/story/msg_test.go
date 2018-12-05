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
	validEvidence := []string{"http://shanesbrain.net"}
	validArgument := "argument"
	msg := NewSubmitStoryMsg(validArgument, validBody, validCategoryID, validCreator, validEvidence, validSource, validStoryType)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Route())
	assert.Equal(t, "submit_story", msg.Type())
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidBodySubmitStoryMsg(t *testing.T) {
	invalidBody := ""
	validCategoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validSource := "http://shanesbrain.net"
	validStoryType := Default
	validEvidence := []string{"http://shanesbrain.net"}
	validArgument := "argument"

	msg := NewSubmitStoryMsg(validArgument, invalidBody, validCategoryID, validCreator, validEvidence, validSource, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(701), err.Code(), err.Error())
}

func TestInValidCreatorSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategoryID := int64(1)
	invalidCreator := sdk.AccAddress([]byte{})
	validSource := "http://shanesbrain.net"
	validStoryType := Default
	validEvidence := []string{"http://shanesbrain.net"}
	validArgument := "argument"

	msg := NewSubmitStoryMsg(validArgument, validBody, validCategoryID, invalidCreator, validEvidence, validSource, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidSourceSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidSource := ""
	validStoryType := Default
	validEvidence := []string{"http://shanesbrain.net"}
	validArgument := ""

	msg := NewSubmitStoryMsg(validArgument, validBody, validCategoryID, validCreator, validEvidence, invalidSource, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(708), err.Code(), err.Error())
}

// ============================================================================

func TestValidAddEvidencetMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validURI := "http://www.truchain.io"
	msg := NewAddEvidenceMsg(validStoryID, validCreator, validURI)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Type())
	assert.Equal(t, "add_evidence", msg.Route())
	assert.Equal(
		t,
		`{"creator":"cosmos1qypq36vzru","story_id":1,"url":"http://www.truchain.io"}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestValidAddArgumentMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validArgument := "another argument"
	msg := NewAddArgumentMsg(validStoryID, validCreator, validArgument)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Type())
	assert.Equal(t, "add_argument", msg.Route())
	assert.Equal(
		t,
		`{"argument":"another argument","creator":"cosmos1qypq36vzru","story_id":1}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestValidFlagStoryMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	msg := NewFlagStoryMsg(validStoryID, validCreator)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Type())
	assert.Equal(t, "flag_story", msg.Route())
	assert.Equal(
		t,
		`{"creator":"cosmos1qypq36vzru","story_id":1}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidStoryIDSubmitEvidencetMsg(t *testing.T) {
	invalidStoryID := int64(-1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validURI := "http://www.truchain.io"
	msg := NewAddEvidenceMsg(invalidStoryID, validCreator, validURI)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(702), err.Code(), err.Error())
}

func TestInValidCreatorSubmitEvidencetMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidCreator := sdk.AccAddress([]byte{})
	validURI := "http://www.truchain.io"
	msg := NewAddEvidenceMsg(validStoryID, invalidCreator, validURI)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidURISubmitEvidencetMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidURI := ""
	msg := NewAddEvidenceMsg(validStoryID, validCreator, invalidURI)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(704), err.Code(), err.Error())
}
