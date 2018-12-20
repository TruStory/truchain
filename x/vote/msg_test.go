package vote

import (
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidCreateVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	creator := sdk.AccAddress([]byte{1, 2})
	evidence := []string{"http://www.trustory.io"}

	msg := NewCreateVoteMsg(storyID, amount, "valid comment", creator, evidence, true)
	err := msg.ValidateBasic()
	assert.Nil(t, err)

	assert.Equal(t, "vote", msg.Route())
	assert.Equal(t, "create_vote", msg.Type())
	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestInValidCreateVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	creator := sdk.AccAddress([]byte{1, 2})
	evidence := []string{"http://www.trustory.io"}

	msg := NewCreateVoteMsg(storyID, amount, "", creator, evidence, true)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, app.ErrInvalidArgumentMsg().Code(), err.Code())
}
