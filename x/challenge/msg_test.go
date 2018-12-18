package challenge

import (
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
	validEvidence := []string{"http://www.cnn.com"}

	msg := NewCreateChallengeMsg(validStoryID, validAmount, validArugment, validCreator, validEvidence)
	err := msg.ValidateBasic()
	assert.Nil(t, err)

	assert.Equal(t, "challenge", msg.Route())
	assert.Equal(t, "create_challenge", msg.Type())
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidStartChallengeMsg(t *testing.T) {
	ctx, _, sk, ck, _ := mockDB()
	validStoryID := createFakeStory(ctx, sk, ck)
	validAmount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	validArugment := "I am against this story because, you know, just cuz."
	validCreator := sdk.AccAddress([]byte{1, 2})
	cnn := "http://www.cnn.com"
	validEvidence := []string{cnn, cnn, cnn, cnn, cnn, cnn, cnn, cnn, cnn, cnn, cnn}

	msg := NewCreateChallengeMsg(validStoryID, validAmount, validArugment, validCreator, validEvidence)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidMsg(msg.Evidence).Code(), err.Code(), "wrong error code")
}
