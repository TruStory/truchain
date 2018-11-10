package vote

import (
	"net/url"
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidCreateVoteMsg(t *testing.T) {
	ctx, _, sk, ck, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	msg := NewCreateVoteMsg(storyID, amount, "valid comment", creator, evidence, true)
	err := msg.ValidateBasic()
	assert.Nil(t, err)

	assert.Equal(t, "vote", msg.Type())
	assert.Equal(t, "create_vote", msg.Route())
	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestInValidCreateVoteMsg(t *testing.T) {
	ctx, _, sk, ck, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	msg := NewCreateVoteMsg(storyID, amount, "", creator, evidence, true)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, app.ErrInvalidCommentMsg().Code(), err.Code())
}
