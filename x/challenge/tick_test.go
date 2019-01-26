package challenge

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, k := fakePendingGameQueue()

	tags := k.NewResponseEndBlock(ctx)
	assert.Equal(t, sdk.Tags{}, tags)
}

func Test_filterExpiredGames(t *testing.T) {
	ctx, k := fakePendingGameQueue()

	q := k.pendingGameList(ctx)
	assert.Equal(t, uint64(1), q.Len())

	err := k.filterExpiredGames(ctx, q)
	assert.Nil(t, err)
}

func Test_removingExpiredGameFromPendingGameQueue(t *testing.T) {
	ctx, k := fakePendingGameQueue()
	q := k.pendingGameList(ctx)

	assert.Equal(t, uint64(1), q.Len())

	// modify challenged expired time in game
	game, _ := k.gameKeeper.Game(ctx, 1)
	expireTime := game.ChallengeExpireTime.AddDate(0, -1, 0)
	game.ChallengeExpireTime = expireTime
	k.gameKeeper.Update(ctx, game)

	err := k.filterExpiredGames(ctx, q)
	assert.Nil(t, err)

	// NOTE: deleting an item from a Cosmos sdk.List type
	// DOES NOT change the result of Len(). So we cannot test
	// if the length of a queue went down after an item is removed.
	// In this case, it is still 1.
	assert.Equal(t, uint64(1), q.Len())

	// but if we iterate the list, we should find nothing..
	q.Iterate(game.ID, func(index uint64) bool {
		assert.Fail(t, "should not find any games")
		return false
	})
}
