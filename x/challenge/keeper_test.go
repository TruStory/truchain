package challenge

import (
	"fmt"
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	ctx, k, _, _, _ := mockDB()

	challenge := Challenge{
		ID:      k.GetNextID(ctx),
		StoryID: int64(5),
	}

	bz := k.GetCodec().MustMarshalBinary(challenge)
	assert.NotNil(t, bz)

	var value Challenge
	k.GetCodec().MustUnmarshalBinary(bz, &value)
	assert.IsType(t, Challenge{}, value, "should be right type")
	assert.Equal(t, challenge.StoryID, value.StoryID, "should be equal")
}

func TestValidKeys(t *testing.T) {
	_, k, _, _, _ := mockDB()

	key := k.GetIDKey(5)
	assert.Equal(t, "challenges:id:5", fmt.Sprintf("%s", key), "should be equal")
}

func TestSetChallenge(t *testing.T) {
	ctx, k, _, _, _ := mockDB()

	challenge := Challenge{ID: int64(5)}
	k.set(ctx, challenge)

	savedChallenge, err := k.Get(ctx, int64(5))
	assert.Nil(t, err)
	assert.Equal(t, challenge.ID, savedChallenge.ID, "should be equal")
}

func TestNewGetChallenge(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	id, err := k.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	challenge, err := k.Get(ctx, id)
	assert.Nil(t, err)

	assert.Equal(t, argument, challenge.Argument, "should match")
}

func TestNewChallenge_Duplicate(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(50))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	challengeAmount, _ := sdk.ParseCoin("10trudex")

	_, err := k.Create(ctx, storyID, challengeAmount, argument, creator, evidence)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, challengeAmount, argument, creator, evidence)
	assert.NotNil(t, err)
	assert.Equal(t, ErrStoryAlreadyChallenged(5).Code(), err.Code())
}

func TestNewChallenge_MultipleChallengers(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(50))
	argument := "test argument is long enough"
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{3, 4})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	challengeAmount, _ := sdk.ParseCoin("10trudex")

	id, err := k.Create(ctx, storyID, challengeAmount, argument, creator1, evidence)
	assert.Nil(t, err)

	challenge, _ := k.Get(ctx, id)

	_, err = k.Update(ctx, challenge.ID, creator2, amount)
	assert.Nil(t, err)
	assert.False(t, bankKeeper.HasCoins(ctx, creator2, sdk.Coins{amount}))

	challenge, _ = k.Get(ctx, id)
	assert.True(t, challenge.Pool.IsEqual(challengeAmount.Plus(amount)))
	assert.True(t, challenge.Started)
}

func TestNewChallenge_ErrIncorrectCategoryCoin(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(15))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	_, err := k.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.NotNil(t, err)
}
