package challenge

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidStartChallengeMsg(t *testing.T) {
	ctx, _, sk, _, _ := mockDB()
	validStoryID := createFakeStory(ctx, sk)
	validAmount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	validArugment := "I am against this story because, you know, just cuz."
	validCreator := sdk.AccAddress([]byte{1, 2})

	msg := NewCreateChallengeMsg(validStoryID, validAmount, validArugment, validCreator)
	err := msg.ValidateBasic()
	assert.Nil(t, err)

	assert.Equal(t, "challenge", msg.Route())
	assert.Equal(t, "create_challenge", msg.Type())
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}
