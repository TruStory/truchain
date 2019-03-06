package vote

import (
	"testing"

	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidCreateVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	creator := sdk.AccAddress([]byte{1, 2})

	msg := NewCreateVoteMsg(storyID, amount, "valid comment", creator, true)
	err := msg.ValidateBasic()
	assert.Nil(t, err)

	assert.Equal(t, "vote", msg.Route())
	assert.Equal(t, "create_vote", msg.Type())
	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}
