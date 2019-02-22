package challenge

import (
	"fmt"
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
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

	id, err := k.Create(ctx, storyID, amount, argument, creator)
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

	id, err := k.Create(ctx, storyID, amount, argument, creator)
	assert.Nil(t, err)

	challenge, err := k.Challenge(ctx, id)
	assert.Nil(t, err)

	assert.Equal(t, amount, challenge.Amount())
}

func TestChallengesByGameID(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	creator2 := sdk.AccAddress([]byte{3, 4})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, argument, creator)
	k.Create(ctx, storyID, amount, argument, creator2)

	story, _ := sk.Story(ctx, storyID)
	challenges, _ := k.ChallengesByStoryID(ctx, story.ID)
	assert.Equal(t, 2, len(challenges))
}

func TestChallengesByStoryIDAndCreator(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, argument, creator)

	challenge, _ := k.ChallengeByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), challenge.ID())
}

func TestTally(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	creator2 := sdk.AccAddress([]byte{3, 4})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, argument, creator)
	k.Create(ctx, storyID, amount, argument, creator2)

	falseVotes, _ := k.Tally(ctx, storyID)

	assert.Equal(t, 2, len(falseVotes))
}

func TestNewChallenge_Duplicate(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(50000000000))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	_, err := k.Create(ctx, storyID, amount, argument, creator)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, argument, creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrDuplicateChallenge(5, creator).Code(), err.Code())
}

func TestNewChallenge_ErrIncorrectCategoryCoin(t *testing.T) {
	ctx, k, sk, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(15000000000))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})

	_, err := k.Create(ctx, storyID, amount, argument, creator)
	assert.NotNil(t, err)
}

func Test_checkThreshold(t *testing.T) {
	ctx, k, storyKeeper, backingKeeper, bankKeeper := mockDB()

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now()})
	storyID := createFakeStory(ctx, storyKeeper)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(10000000000))
	argument := "test argument right here"
	backer1 := fakeFundedCreator(ctx, bankKeeper)
	backer2 := fakeFundedCreator(ctx, bankKeeper)
	challenger1 := fakeFundedCreator(ctx, bankKeeper)
	challenger2 := fakeFundedCreator(ctx, bankKeeper)
	duration := 5 * 24 * time.Hour

	_, err := backingKeeper.Create(
		ctx, storyID, amount, argument, backer1, duration)
	assert.Nil(t, err)

	_, err = backingKeeper.Create(
		ctx, storyID, amount, argument, backer2, duration)
	assert.Nil(t, err)

	_, err = k.Create(
		ctx, storyID, amount, argument, challenger1)
	assert.Nil(t, err)

	_, err = k.Create(
		ctx, storyID, amount, argument, challenger2)
	assert.Nil(t, err)

	err = k.checkThreshold(ctx, storyID)
	assert.Nil(t, err)

	story, _ := storyKeeper.Story(ctx, storyID)
	assert.Equal(t, story.State.String(), "Voting")
}
