package challenge

import (
	"fmt"
	"testing"

	params "github.com/TruStory/truchain/parameters"
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
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
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
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
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
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
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
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	k.Create(ctx, storyID, amount, argument, creator)

	challenge, _ := k.ChallengeByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), challenge.ID())
}

func TestTally(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
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
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(50000000000))
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

// TODO [shanev]: Add this to game in https://github.com/TruStory/truchain/issues/387
// func TestNewChallenge_MultipleChallengers(t *testing.T) {
// 	ctx, k, sk, ck, bankKeeper := mockDB()

// 	storyID := createFakeStory(ctx, sk, ck)
// 	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(50000000000))
// 	argument := "test argument is long enough"
// 	creator1 := sdk.AccAddress([]byte{1, 2})
// 	creator2 := sdk.AccAddress([]byte{3, 4})

// 	// give user some funds
// 	bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount})
// 	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

// 	id, err := k.Create(ctx, storyID, amount, argument, creator1)
// 	assert.Nil(t, err)

// 	challenge, _ := k.Challenge(ctx, id)

// 	_, err = k.Create(ctx, challenge.ID(), amount, argument, creator2)
// 	assert.Nil(t, err)
// 	assert.False(t, bankKeeper.HasCoins(ctx, creator2, sdk.Coins{amount}))

// 	// check game pool amount
// 	story, _ := k.storyKeeper.Story(ctx, storyID)
// 	game, _ := k.gameKeeper.Game(ctx, story.GameID)

// 	assert.True(t, game.ChallengePool.IsEqual(amount.Plus(amount)))
// }

func TestNewChallenge_ErrIncorrectCategoryCoin(t *testing.T) {
	ctx, k, sk, ck, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(15000000000))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})

	_, err := k.Create(ctx, storyID, amount, argument, creator)
	assert.NotNil(t, err)
}
