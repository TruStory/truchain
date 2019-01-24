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

func Test_pendingGameQueue(t *testing.T) {
	ctx, k := fakePendingGameQueue()

	q := k.pendingGameQueue(ctx)
	assert.Equal(t, uint64(1), q.List.Len())

	err := k.checkPendingQueue(ctx, q)
	assert.Nil(t, err)
}

func Test_removingExpiredGameFromPendingGameQueue(t *testing.T) {
	ctx, k := fakePendingGameQueue()
	q := k.pendingGameQueue(ctx)

	assert.Equal(t, uint64(1), q.List.Len())

	// modify challenged expired time in game
	game, _ := k.gameKeeper.Game(ctx, 1)
	expireTime := game.ChallengeExpireTime.AddDate(0, -1, 0)
	game.ChallengeExpireTime = expireTime
	k.gameKeeper.Update(ctx, game)

	err := k.checkPendingQueue(ctx, q)
	assert.Nil(t, err)

	// NOTE: popping/deleting an item from a Cosmos queue
	// DOES NOT change the result of Len(). So we cannot test
	// if the length of a queue went down after an item is removed.
	// assert.Equal(t, uint64(0), q.List.Len())
}
