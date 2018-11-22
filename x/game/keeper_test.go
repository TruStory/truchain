package game

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {
	ctx, k, storyKeeper, categoryKeeper, _ := mockDB()

	storyID := createFakeStory(ctx, storyKeeper, categoryKeeper)

	creator := sdk.AccAddress([]byte{1, 2})

	gameID, err := k.Create(ctx, storyID, creator)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), gameID)
}

func TestRegisterChallenge(t *testing.T) {
	ctx, k, storyKeeper, categoryKeeper, _ := mockDB()

	storyID := createFakeStory(ctx, storyKeeper, categoryKeeper)
	creator := sdk.AccAddress([]byte{1, 2})
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")
	err := k.RegisterChallenge(ctx, gameID, amount)
	assert.Nil(t, err)

	game, _ := k.Game(ctx, gameID)
	assert.Equal(t, sdk.NewInt(50), game.ChallengeThreshold.Amount)
}

func TestRegisterVoteGameNotStarted(t *testing.T) {
	ctx, k, storyKeeper, categoryKeeper, _ := mockDB()

	storyID := createFakeStory(ctx, storyKeeper, categoryKeeper)
	creator := sdk.AccAddress([]byte{1, 2})
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")
	k.RegisterChallenge(ctx, gameID, amount)

	quorum := DefaultParams().VoterQuorum - 1
	for i := 0; i < int(quorum); i = i + 1 {
		k.RegisterVote(ctx, gameID)
	}

	game, _ := k.Game(ctx, gameID)
	assert.False(t, game.Started())
}

func TestRegisterVoteGameStarted(t *testing.T) {
	ctx, k, storyKeeper, categoryKeeper, _ := mockDB()

	storyID := createFakeStory(ctx, storyKeeper, categoryKeeper)
	creator := sdk.AccAddress([]byte{1, 2})
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")
	k.RegisterChallenge(ctx, gameID, amount)

	quorum := DefaultParams().VoterQuorum + 1
	for i := 0; i < int(quorum); i = i + 1 {
		k.RegisterVote(ctx, gameID)
	}

	game, _ := k.Game(ctx, gameID)
	assert.True(t, game.Started())
}

func TestSetGame(t *testing.T) {
	ctx, k, _, _, _ := mockDB()

	game := Game{ID: int64(5)}
	k.set(ctx, game)

	savedGame, err := k.Game(ctx, int64(5))
	assert.Nil(t, err)
	assert.Equal(t, game.ID, savedGame.ID)
}
