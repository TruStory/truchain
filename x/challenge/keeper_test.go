package challenge

import (
	"fmt"
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidKeys(t *testing.T) {
	_, k, _, _, _ := mockDB()

	key := k.GetIDKey(5)
	assert.Equal(t, "challenges:id:5", fmt.Sprintf("%s", key), "should be equal")
}

func TestNewGetChallenge(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	id, err := k.Create(ctx, storyID, amount, 0, argument, creator)
	assert.Nil(t, err)

	challenge, err := k.Challenge(ctx, id)
	assert.Nil(t, err)

	assert.Equal(t, amount, challenge.Amount())
}

func TestNewGetChallengeUsingTruStake(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	id, err := k.Create(ctx, storyID, amount, 0, argument, creator)
	assert.Nil(t, err)

	challenge, err := k.Challenge(ctx, id)
	assert.Nil(t, err)

	assert.Equal(t, amount, challenge.Amount())
}

func TestChallengesByStoryID(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	creator2 := sdk.AccAddress([]byte{3, 4})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, 0, argument, creator)

	story, _ := sk.Story(ctx, storyID)
	challenges, _ := k.ChallengesByStoryID(ctx, story.ID)
	assert.Equal(t, 1, len(challenges))
}

func TestChallengesByStoryIDAndCreator(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, 0, argument, creator)

	challenge, _ := k.ChallengeByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), challenge.ID())
}

func TestNewChallenge_Duplicate(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	_, err := k.Create(ctx, storyID, amount, 0, argument, creator)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, 0, argument, creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrDuplicateChallenge(5, creator).Code(), err.Code())
}

func TestNewChallenge_ErrIncorrectCategoryCoin(t *testing.T) {
	ctx, k, sk, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(15000000000))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})

	_, err := k.Create(ctx, storyID, amount, 0, argument, creator)
	assert.NotNil(t, err)
}

func Test_ChallengersByStoryID(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{1, 2, 3, 4})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	_, err := k.Create(ctx, storyID, amount, 0, argument, creator)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, 0, argument, creator2)
	assert.Nil(t, err)

	challengers, err := k.ChallengersByStoryID(ctx, storyID)
	assert.Nil(t, err)
	assert.Subset(t, []sdk.Address{creator, creator2}, challengers)
}
