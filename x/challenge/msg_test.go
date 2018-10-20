package challenge

import (
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidStartChallengeMsg(t *testing.T) {
	ctx, _, sk, ck, _ := mockDB()
	validStoryID := createFakeStory(ctx, sk, ck)
	validAmount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	validArugment := "I am against this story because, you know, just cuz."
	validCreator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	validEvidence := []url.URL{*cnn}
	validReason := False

	msg := NewStartChallengeMsg(validStoryID, validAmount, validArugment, validCreator, validEvidence, validReason)
	err := msg.ValidateBasic()
	assert.Nil(t, err)

	assert.Equal(t, "challenge", msg.Type())
	assert.Equal(t, "start_challenge", msg.Name())
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidStartChallengeMsg(t *testing.T) {
	ctx, _, sk, ck, _ := mockDB()
	validStoryID := createFakeStory(ctx, sk, ck)
	validAmount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	validArugment := "I am against this story because, you know, just cuz."
	validCreator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	validEvidence := []url.URL{*cnn, *cnn, *cnn, *cnn, *cnn, *cnn, *cnn, *cnn, *cnn, *cnn, *cnn}
	validReason := False

	msg := NewStartChallengeMsg(validStoryID, validAmount, validArugment, validCreator, validEvidence, validReason)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidMsg(msg.Evidence).Code(), err.Code(), "wrong error code")
}

func TestValidUpdateChallengeMsg(t *testing.T) {
	ctx, _, sk, ck, _ := mockDB()
	validStoryID := createFakeStory(ctx, sk, ck)
	validAmount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	validArugment := "I am against this story because, you know, just cuz."
	validCreator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	validEvidence := []url.URL{*cnn}
	validReason := False

	msg := NewStartChallengeMsg(validStoryID, validAmount, validArugment, validCreator, validEvidence, validReason)
	err := msg.ValidateBasic()
	assert.Nil(t, err)

	updateMsg := NewUpdateChallengeMsg(1, validAmount, validCreator)
	err = updateMsg.ValidateBasic()
	assert.Nil(t, err)
}
