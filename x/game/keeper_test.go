package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGame(t *testing.T) {
	ctx, k, _, _, _ := mockDB()

	game := Game{ID: int64(5)}
	k.set(ctx, game)

	savedGame, err := k.Get(ctx, int64(5))
	assert.Nil(t, err)
	assert.Equal(t, game.ID, savedGame.ID)
}
