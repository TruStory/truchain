package game

import (
	"net/url"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var creator = sdk.AccAddress([]byte{1, 2})

func evidence() []url.URL {
	testURL, _ := url.Parse("http://www.trustory.io")
	return []url.URL{*testURL}
}

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
	err := k.RegisterChallenge(ctx, gameID, amount)
	assert.Nil(t, err)

	game, _ := k.Game(ctx, gameID)
	assert.Equal(t, sdk.NewInt(50), game.ChallengePool.Amount)
	assert.Equal(t, true, game.Started)
}

func TestRegisterChallengeNoBackersNotMeetMinChallenge(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("5trudex")
	err := k.RegisterChallenge(ctx, gameID, amount)
	assert.Nil(t, err)

	game, _ := k.Game(ctx, gameID)
	assert.Equal(t, sdk.NewInt(5), game.ChallengePool.Amount)
	assert.Equal(t, false, game.Started)
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
	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration, evidence())

	// challenge with 33trudex (33% of total backings)
	amount, _ = sdk.ParseCoin("33trudex")
	err := k.RegisterChallenge(ctx, gameID, amount)
	assert.Nil(t, err)

	game, _ := k.Game(ctx, gameID)
	assert.Equal(t, true, game.Started)
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
	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration, evidence())

	// challenge with 32trudex (32% of total backings)
	amount, _ = sdk.ParseCoin("32trudex")
	err := k.RegisterChallenge(ctx, gameID, amount)
	assert.Nil(t, err)

	game, _ := k.Game(ctx, gameID)
	assert.Equal(t, false, game.Started)
}

func TestRegisterVoteGameStarted(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")
	k.RegisterChallenge(ctx, gameID, amount)

	quorum := DefaultParams().VoteQuorum - 1
	for i := 0; i < int(quorum); i = i + 1 {
		k.RegisterVote(ctx, gameID)
	}

	game, _ := k.Game(ctx, gameID)
	assert.True(t, game.Started)
}

func TestRegisterVoteGameEnded(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")
	k.RegisterChallenge(ctx, gameID, amount)

	quorum := DefaultParams().VoteQuorum + 1
	for i := 0; i < int(quorum); i = i + 1 {
		k.RegisterVote(ctx, gameID)
	}

	game, _ := k.Game(ctx, gameID)
	assert.True(t, game.Started)
	endTime := game.EndTime.Add(20 * 24 * time.Hour)
	assert.True(t, game.Ended(endTime))
}

func TestRegisterVoteGameExpired(t *testing.T) {
	ctx, k, categoryKeeper := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")
	k.RegisterChallenge(ctx, gameID, amount)

	quorum := DefaultParams().VoteQuorum - 1
	for i := 0; i < int(quorum); i = i + 1 {
		k.RegisterVote(ctx, gameID)
	}

	game, _ := k.Game(ctx, gameID)
	endTime := game.EndTime.Add(20 * 24 * time.Hour)
	assert.True(t, game.Expired(endTime))
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

func Test_challengeThresholdWithBacking(t *testing.T) {
	_, k, _ := mockDB()
	amt := k.ChallengeThreshold(sdk.NewCoin("trudex", sdk.NewInt(100)))

	assert.Equal(t, "33trudex", amt.String())
}
