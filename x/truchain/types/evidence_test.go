package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidSubmitEvidencetMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validURI := "http://www.truchain.io"
	msg := NewSubmitEvidenceMsg(validStoryID, validCreator, validURI)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "SubmitEvidence", msg.Type())
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

	assert.Equal(t, sdk.CodeType(703), err.Code(), err.Error())
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

	assert.Equal(t, sdk.CodeType(707), err.Code(), err.Error())
}
