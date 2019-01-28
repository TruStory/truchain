package game

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var creator = sdk.AccAddress([]byte{1, 2})

func TestCreateGame(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)

	gameID, err := k.Create(ctx, storyID, creator)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), gameID)
}

func TestRegisterChallengeNoBackersMeetMinChallenge(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")
	err := k.AddToChallengePool(ctx, gameID, amount)
	assert.Nil(t, err)

	game, _ := k.Game(ctx, gameID)
	assert.Equal(t, sdk.NewInt(50), game.ChallengePool.Amount)
}

func TestRegisterChallengeNoBackersNotMeetMinChallenge(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("5trudex")
	err := k.AddToChallengePool(ctx, gameID, amount)
	assert.Nil(t, err)

	game, _ := k.Game(ctx, gameID)
	assert.Equal(t, sdk.NewInt(5), game.ChallengePool.Amount)
}

func TestRegisterChallengeHaveBackersMeetThreshold(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)
	amount, _ := sdk.ParseCoin("100trudex")
	argument := "cool story brew"

	// back story with 100trudex
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	duration := 30 * 24 * time.Hour
	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration)

	// challenge with 33trudex (33% of total backings)
	amount, _ = sdk.ParseCoin("33trudex")
	err := k.AddToChallengePool(ctx, gameID, amount)
	assert.Nil(t, err)
}

func TestRegisterChallengeHaveBackersNotMeetThreshold(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)
	amount, _ := sdk.ParseCoin("100trudex")
	argument := "cool story brew"

	// back story with 100trudex
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	duration := 30 * 24 * time.Hour
	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration)

	// challenge with 32trudex (32% of total backings)
	amount, _ = sdk.ParseCoin("32trudex")
	err := k.AddToChallengePool(ctx, gameID, amount)
	assert.Nil(t, err)
}

func TestSetGame(t *testing.T) {
	ctx, k, _ := mockDB()

	game := Game{ID: int64(5)}
	k.set(ctx, game)

	savedGame, err := k.Game(ctx, int64(5))
	assert.Nil(t, err)
	assert.Equal(t, game.ID, savedGame.ID)
}

func Test_challengeThresholdNoBacking(t *testing.T) {
	_, k, _ := mockDB()
	amt := k.ChallengeThreshold(sdk.NewCoin("trudex", sdk.ZeroInt()))

	assert.Equal(t, "10trudex", amt.String())
}

// challenge threshold should not go below min challenge stake
// if challenge threshold is 1/3 of backing and min challenge stake is 10
// when a story has backing amount of 21, challenge threshold should be 10, not 7
// then instead
func Test_challengeThresholdWithSmallBacking(t *testing.T) {
	_, k, _ := mockDB()
	amt := k.ChallengeThreshold(sdk.NewCoin("trudex", sdk.NewInt(21)))

	assert.Equal(t, "10trudex", amt.String())
}

func Test_challengeThresholdWithBacking(t *testing.T) {
	_, k, _ := mockDB()
	amt := k.ChallengeThreshold(sdk.NewCoin("trudex", sdk.NewInt(100)))

	assert.Equal(t, "33trudex", amt.String())
}

func Test_start(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)
	amount, _ := sdk.ParseCoin("100trudex")
	argument := "cool story brew"

	// back story with 100trudex
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	duration := 30 * 24 * time.Hour
	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration)

	// challenge with 33trudex (33% of total backings)
	amount, _ = sdk.ParseCoin("33trudex")
	err := k.AddToChallengePool(ctx, gameID, amount)
	assert.Nil(t, err)

	// test queue sizes
	assert.Equal(t, uint64(1), k.pendingList(ctx).Len())
	assert.Equal(t, uint64(1), k.queue(ctx).List.Len())
}
