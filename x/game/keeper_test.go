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

func TestUpdateGame(t *testing.T) {
	ctx, k, storyKeeper, categoryKeeper, _ := mockDB()

	storyID := createFakeStory(ctx, storyKeeper, categoryKeeper)
	creator := sdk.AccAddress([]byte{1, 2})
	gameID, _ := k.Create(ctx, storyID, creator)

	amount, _ := sdk.ParseCoin("50trudex")

	gameID, _ = k.Update(ctx, gameID, amount)

	game, _ := k.Get(ctx, gameID)
	assert.Equal(t, sdk.NewInt(50), game.Pool.Amount)
}

func TestSetGame(t *testing.T) {
	ctx, k, _, _, _ := mockDB()

	game := Game{ID: int64(5)}
	k.Set(ctx, game)

	savedGame, err := k.Get(ctx, int64(5))
	assert.Nil(t, err)
	assert.Equal(t, game.ID, savedGame.ID)
}
