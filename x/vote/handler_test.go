package vote

import (
	"encoding/json"
	"testing"

	"github.com/TruStory/truchain/x/story"

	"github.com/TruStory/truchain/types"
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	msg := NewCreateVoteMsg(storyID, amount, "valid comment", creator, true)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID)
}

func Test_InvalidCreateVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	// Amout is zero
	msg := NewCreateVoteMsg(storyID, sdk.NewCoin(app.StakeDenom, sdk.NewInt(0)), "valid comment", creator, true)
	assert.NotNil(t, msg)
	res := h(ctx, msg)
	assert.Equal(t, sdk.CodeInsufficientFunds, res.Code)

	// Creator is invalid
	msg = NewCreateVoteMsg(storyID, amount, "valid comment", sdk.AccAddress([]byte{}), true)
	assert.NotNil(t, msg)
	res = h(ctx, msg)
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)

	// StoryID is invalid
	msg = NewCreateVoteMsg(0, amount, "valid comment", sdk.AccAddress([]byte{}), true)
	assert.NotNil(t, msg)
	res = h(ctx, msg)
	assert.Equal(t, story.CodeInvalidStoryID, res.Code)

	// Err
	storyID = createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	msg = NewCreateVoteMsg(storyID, amount, "valid comment", creator, true)
	assert.NotNil(t, msg)
	res = h(ctx, msg)
	assert.Equal(t, CodeVotingNotStarted, res.Code)
}

func TestToggleVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	msg := NewToggleVoteMsg(storyID, amount, "valid comment", creator, true)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID)
}

func Test_InvalidToggleVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	// Creator is invalid
	msg := NewToggleVoteMsg(storyID, amount, "valid comment", sdk.AccAddress([]byte{}), true)
	assert.NotNil(t, msg)
	res := h(ctx, msg)
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)

	// StoryID is invalid
	msg = NewToggleVoteMsg(0, amount, "valid comment", sdk.AccAddress([]byte{}), true)
	assert.NotNil(t, msg)
	res = h(ctx, msg)
	assert.Equal(t, story.CodeInvalidStoryID, res.Code)

	// Err
	storyID = createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	msg = NewToggleVoteMsg(storyID, amount, "valid comment", creator, true)
	assert.NotNil(t, msg)
	res = h(ctx, msg)
	assert.Equal(t, CodeInvalidStoryState, res.Code)
}
