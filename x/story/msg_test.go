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
	validStoryType := Default
	msg := NewSubmitStoryMsg(validBody, validCategoryID, validCreator, validStoryType)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Route())
	assert.Equal(t, "submit_story", msg.Type())
	assert.Equal(
		t,
		`{"body":"This is a valid story body @shanev amirite?","category_id":1,"creator":"cosmos1qypq36vzru","story_type":0}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidBodySubmitStoryMsg(t *testing.T) {
	invalidBody := ""
	validCategoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validStoryType := Default
	msg := NewSubmitStoryMsg(invalidBody, validCategoryID, validCreator, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(701), err.Code(), err.Error())
}

func TestInValidCreatorSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategoryID := int64(1)
	invalidCreator := sdk.AccAddress([]byte{})
	validStoryType := Default
	msg := NewSubmitStoryMsg(validBody, validCategoryID, invalidCreator, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

// ============================================================================

func TestValidSubmitEvidencetMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validURI := "http://www.truchain.io"
	msg := NewSubmitEvidenceMsg(validStoryID, validCreator, validURI)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "story", msg.Type())
	assert.Equal(t, "submit_evidence", msg.Route())
	assert.Equal(
		t,
		`{"creator":"cosmos1qypq36vzru","story_id":1,"url":"http://www.truchain.io"}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidStoryIDSubmitEvidencetMsg(t *testing.T) {
	invalidStoryID := int64(-1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validURI := "http://www.truchain.io"
	msg := NewSubmitEvidenceMsg(invalidStoryID, validCreator, validURI)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(702), err.Code(), err.Error())
}

func TestInValidCreatorSubmitEvidencetMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidCreator := sdk.AccAddress([]byte{})
	validURI := "http://www.truchain.io"
	msg := NewSubmitEvidenceMsg(validStoryID, invalidCreator, validURI)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidURISubmitEvidencetMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidURI := ""
	msg := NewSubmitEvidenceMsg(validStoryID, validCreator, invalidURI)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(704), err.Code(), err.Error())
}
