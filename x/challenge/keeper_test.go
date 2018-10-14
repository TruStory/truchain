package challenge

import (
	"fmt"
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	ctx, k, _, _ := mockDB()

	challenge := Challenge{
		ID:      k.GetNextID(ctx, k.storeKey),
		StoryID: int64(5),
	}

	bz := k.marshal(challenge)
	assert.NotNil(t, bz)

	value := k.unmarshal(bz)
	assert.IsType(t, Challenge{}, value, "should be right type")
	assert.Equal(t, challenge.StoryID, value.StoryID, "should be equal")
}

func TestValidKeys(t *testing.T) {
	_, k, _, _ := mockDB()

	key := getChallengeIDKey(k, 5)
	assert.Equal(t, "challenges:id:5", fmt.Sprintf("%s", key), "should be equal")
}

func TestSetChallenge(t *testing.T) {
	ctx, k, _, _ := mockDB()

	challenge := Challenge{ID: int64(5)}
	k.setChallenge(ctx, challenge)

	savedChallenge, err := k.GetChallenge(ctx, int64(5))
	assert.Nil(t, err)
	assert.Equal(t, challenge.ID, savedChallenge.ID, "should be equal")
}

func TestNewGetChallenge(t *testing.T) {
	ctx, k, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	id, err := k.NewChallenge(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	challenge, _ := k.GetChallenge(ctx, id)
	assert.Equal(t, argument, challenge.Arugment, "should match")
}
