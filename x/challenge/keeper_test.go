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

	id, err := k.Create(ctx, storyID, amount, 0, argument, creator, false)
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

	id, err := k.Create(ctx, storyID, amount, 0, argument, creator, false)
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

	k.Create(ctx, storyID, amount, 0, argument, creator, false)

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

	k.Create(ctx, storyID, amount, 0, argument, creator, false)

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

	k.Create(ctx, storyID, amount, 0, argument, creator, false)
	falseVotes, _ := k.Tally(ctx, storyID)

	assert.Equal(t, 1, len(falseVotes))
}

func TestNewChallenge_Duplicate(t *testing.T) {
	ctx, k, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	_, err := k.Create(ctx, storyID, amount, 0, argument, creator, false)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, 0, argument, creator, false)
	assert.NotNil(t, err)
	assert.Equal(t, ErrDuplicateChallenge(5, creator).Code(), err.Code())
}

func TestNewChallenge_ErrIncorrectCategoryCoin(t *testing.T) {
	ctx, k, sk, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(15000000000))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})

	_, err := k.Create(ctx, storyID, amount, 0, argument, creator, false)
	assert.NotNil(t, err)
}

func Test_ChallengeDelete(t *testing.T) {
	ctx, k, storyKeeper, _, bankKeeper := mockDB()
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now()})
	storyID := createFakeStory(ctx, storyKeeper)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(10000000000))
	argument := "test argument right here"
	challenger1 := fakeFundedCreator(ctx, bankKeeper)
	totalCoins := bankKeeper.GetCoins(ctx, challenger1)

	id, err := k.Create(
		ctx, storyID, amount, 0, argument, challenger1, false)
	assert.NoError(t, err)

	challenge, err := k.Challenge(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, challenge.Vote)
	assert.Equal(t, totalCoins.Minus(sdk.Coins{amount}), bankKeeper.GetCoins(ctx, challenger1), "coins should have been deducted")

	k.Delete(ctx, challenge)

	_, err = k.Challenge(ctx, id)
	assert.Equal(t, CodeNotFound, err.Code())
	assert.Equal(t, totalCoins, bankKeeper.GetCoins(ctx, challenger1), "coins should have been added back")
}
