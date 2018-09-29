package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

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
		`{"amount":{"amount":"100","denom":"trusomecoin"},"creator":"cosmos1qypq36vzru","story_id":1,"vote":true}`,
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
