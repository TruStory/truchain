package challenge

import (
	"encoding/binary"
	"encoding/json"
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestSubmitChallengeMsg(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument"
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	msg := NewCreateChallengeMsg(storyID, amount, argument, creator)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	result := &app.StakeNotificationResult{}
	err := json.Unmarshal(res.Data, result)
	assert.NotNil(t, msg)

	story, err := sk.Story(ctx, storyID)
	assert.NoError(t, err)

	expected := &app.StakeNotificationResult{
		MsgResult: app.MsgResult{ID: int64(1)},
		Amount:    amount,
		StoryID:   storyID,
		From:      creator,
		To:        story.Creator,
	}

	assert.Equal(t, expected, result)
}

func TestLikeChallengeMsg(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument"
	challengeCreator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	// give user some funds
	bankKeeper.AddCoins(ctx, challengeCreator, sdk.Coins{amount})

	msg := NewCreateChallengeMsg(storyID, amount, argument, challengeCreator)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	result := &app.StakeNotificationResult{}
	_ = json.Unmarshal(res.Data, result)

	story, err := sk.Story(ctx, storyID)
	assert.NoError(t, err)

	expected := &app.StakeNotificationResult{
		MsgResult: app.MsgResult{ID: int64(1)},
		Amount:    amount,
		StoryID:   storyID,
		From:      challengeCreator,
		To:        story.Creator,
	}

	assert.Equal(t, expected, result)

	challenge, err := k.Challenge(ctx, result.ID)
	assert.NoError(t, err)
	likeCreator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	// give user some funds
	bankKeeper.AddCoins(ctx, likeCreator, sdk.Coins{amount})

	// Test Like a challenge.
	likeMsg := NewLikeChallengeArgumentMsg(challenge.ArgumentID, likeCreator, amount)

	res = h(ctx, likeMsg)

	likeResult := &app.StakeNotificationResult{}
	_ = json.Unmarshal(res.Data, likeResult)

	stakeToCredRatio := k.stakeKeeper.GetParams(ctx).StakeToCredRatio
	expectedCred := sdk.NewCoin("crypto", amount.Amount.Quo(stakeToCredRatio))

	expectedLikeResult := &app.StakeNotificationResult{
		MsgResult: app.MsgResult{ID: int64(2)},
		Amount:    amount,
		StoryID:   storyID,
		From:      likeCreator,
		To:        challengeCreator,
		Cred:      &expectedCred,
	}

	assert.Equal(t, expectedLikeResult, likeResult)
}

func TestSubmitChallengeMsg_ErrInsufficientFunds(t *testing.T) {
	ctx, k, sk, _, _ := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})

	msg := NewCreateChallengeMsg(storyID, amount, argument, creator)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(0), x, "incorrect result data")
}

func TestSubmitChallengeMsg_ErrInsufficientChallengeAmount(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin("trudex", sdk.NewInt(1))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	msg := NewCreateChallengeMsg(storyID, amount, argument, creator)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(0), x, "incorrect result data")
}
