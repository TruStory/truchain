package challenge

import (
	"fmt"
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidKeys(t *testing.T) {
	_, k, _, _, _ := mockDB()

	key := k.GetIDKey(5)
	assert.Equal(t, "challenges:id:5", fmt.Sprintf("%s", key), "should be equal")
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

	challenge, err := k.Challenge(ctx, id)
	assert.Nil(t, err)

	assert.Equal(t, amount, challenge.Amount())
}

func TestNewGetChallengeUsingTruStake(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(15))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	id, err := k.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	challenge, err := k.Challenge(ctx, id)
	assert.Nil(t, err)

	expectedCoin := sdk.NewCoin("trudex", amount.Amount)
	assert.Equal(t, expectedCoin, challenge.Amount())
}

func TestChallengesByGameID(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	argument := "test argument is long enough"
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	creator2 := sdk.AccAddress([]byte{3, 4})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, argument, creator, evidence)
	k.Create(ctx, storyID, amount, argument, creator2, evidence)

	story, _ := sk.Story(ctx, storyID)
	challenges, _ := k.ChallengesByGameID(ctx, story.GameID)
	assert.Equal(t, 2, len(challenges))
}

func TestChallengesByStoryIDAndCreator(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	argument := "test argument is long enough"
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, argument, creator, evidence)

	challenge, _ := k.ChallengeByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), challenge.ID())
}

func TestTally(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	argument := "test argument is long enough"
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	creator2 := sdk.AccAddress([]byte{3, 4})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, argument, creator, evidence)
	k.Create(ctx, storyID, amount, argument, creator2, evidence)

	falseVotes, _ := k.Tally(ctx, storyID)

	assert.Equal(t, 2, len(falseVotes))
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
	assert.Equal(t, ErrDuplicateChallenge(5, creator).Code(), err.Code())
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

	challenge, _ := k.Challenge(ctx, id)

	_, err = k.Create(ctx, challenge.ID(), amount, argument, creator2, evidence)
	assert.Nil(t, err)
	assert.False(t, bankKeeper.HasCoins(ctx, creator2, sdk.Coins{amount}))

	// check game pool amount
	story, _ := k.storyKeeper.Story(ctx, storyID)
	game, _ := k.gameKeeper.Game(ctx, story.GameID)

	assert.True(t, game.ChallengePool.IsEqual(challengeAmount.Plus(amount)))
	assert.True(t, game.Started)
}

func TestNewChallenge_ErrIncorrectCategoryCoin(t *testing.T) {
	ctx, k, sk, ck, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(15))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	_, err := k.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.NotNil(t, err)
}
