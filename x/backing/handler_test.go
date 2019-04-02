package backing

import (
	"encoding/json"
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestBackStoryMsg_FailInsufficientFunds(t *testing.T) {
	ctx, bk, sk, ck, _, _ := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	msg := NewBackStoryMsg(storyID, amount, argument, creator)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	assert.Equal(t, sdk.CodeInsufficientFunds, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)

}

func TestBackStoryMsg(t *testing.T) {
	ctx, bk, sk, ck, _, am := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := createFakeFundedAccount(ctx, am, sdk.Coins{amount})
	msg := NewBackStoryMsg(storyID, amount, argument, creator)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	result := &app.StakeNotificationResult{}
	_ = json.Unmarshal(res.Data, result)

	expected := &app.StakeNotificationResult{
		MsgResult: app.MsgResult{ID: int64(1), Amount: amount},
		StoryID:   storyID,
		From:      creator,
		To:        sdk.AccAddress([]byte{1, 2}),
	}

	assert.Equal(t, expected, result)
}

func TestByzantineMsg(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()
	h := NewHandler(bk)
	assert.NotNil(t, h)
	res := h(ctx, nil)
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)
}
