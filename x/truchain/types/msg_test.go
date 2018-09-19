package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// ============================================================================

func TestValidPlaceBondMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewPlaceBondMsg(validStoryID, validStake, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "PlaceBond", msg.Type())
	assert.Equal(
		t,
		`{"amount":{"amount":"100","denom":"trusomecoin"},"creator":"cosmosaccaddr1qypq8zs0ka","period":259200000000000,"story_id":1}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInvalidStoryIdPlaceBondMsg(t *testing.T) {
	invalidStoryID := int64(-1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewPlaceBondMsg(invalidStoryID, validStake, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(703), err.Code(), err.Error())
}

func TestInvalidAddressPlaceBondMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	invalidCreator := sdk.AccAddress([]byte{})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewPlaceBondMsg(validStoryID, validStake, invalidCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidStakePlaceBondMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(0)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	validPeriod := time.Duration(3 * 24 * time.Hour)
	msg := NewPlaceBondMsg(validStoryID, invalidStake, validCreator, validPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(705), err.Code(), err.Error())
}

func TestInValidBondPeriodPlaceBondMsg(t *testing.T) {
	validStoryID := int64(1)
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidPeriod := time.Duration(0 * time.Hour)
	msg := NewPlaceBondMsg(validStoryID, validStake, validCreator, invalidPeriod)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(706), err.Code(), err.Error())
}

// ============================================================================

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

// ============================================================================

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
		`{"creator":"cosmosaccaddr1qypq8zs0ka","story_id":1,"url":"http://www.truchain.io"}`,
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

// ============================================================================

func TestValidSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategory := DEX
	validCreator := sdk.AccAddress([]byte{1, 2})
	validEscrow := sdk.AccAddress([]byte{3, 4})
	validStoryType := Default
	msg := NewSubmitStoryMsg(validBody, validCategory, validCreator, validEscrow, validStoryType)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "SubmitStory", msg.Type())
	assert.Equal(
		t,
		`{"body":"This is a valid story body @shanev amirite?","category":3,"creator":"cosmosaccaddr1qypq8zs0ka","escrow":"cosmosaccaddr1qvzqxum5ug","story_type":0}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidBodySubmitStoryMsg(t *testing.T) {
	invalidBody := ""
	validCategory := DEX
	validCreator := sdk.AccAddress([]byte{1, 2})
	validEscrow := sdk.AccAddress([]byte{3, 4})
	validStoryType := Default
	msg := NewSubmitStoryMsg(invalidBody, validCategory, validCreator, validEscrow, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(702), err.Code(), err.Error())
}

func TestInValidCreatorSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategory := DEX
	invalidCreator := sdk.AccAddress([]byte{})
	validEscrow := sdk.AccAddress([]byte{3, 4})
	validStoryType := Default
	msg := NewSubmitStoryMsg(validBody, validCategory, invalidCreator, validEscrow, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

// ============================================================================

func TestValidVoteMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validVote := true
	msg := NewVoteMsg(validStoryID, validCreator, validStake, validVote)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "Vote", msg.Type())
	assert.Equal(
		t,
		`{"amount":{"amount":"100","denom":"trusomecoin"},"creator":"cosmosaccaddr1qypq8zs0ka","story_id":1,"vote":true}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidStoryIDVoteMsg(t *testing.T) {
	invalidStoryID := int64(-1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validVote := true
	msg := NewVoteMsg(invalidStoryID, validCreator, validStake, validVote)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(703), err.Code(), err.Error())
}

func TestInValidAddressVoteMsg(t *testing.T) {
	validStoryID := int64(1)
	invalidCreator := sdk.AccAddress([]byte{})
	validStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(100)}
	validVote := true
	msg := NewVoteMsg(validStoryID, invalidCreator, validStake, validVote)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}

func TestInValidStakeVoteMsg(t *testing.T) {
	validStoryID := int64(1)
	validCreator := sdk.AccAddress([]byte{1, 2})
	invalidStake := sdk.Coin{Denom: "trusomecoin", Amount: sdk.NewInt(0)}
	validVote := true
	msg := NewVoteMsg(validStoryID, validCreator, invalidStake, validVote)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(705), err.Code(), err.Error())
}
