package vote

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	tags := k.NewResponseEndBlock(ctx)
	assert.Equal(t, sdk.Tags{}, tags)
}

func Test_checkGames(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	qStore := ctx.KVStore(k.activeGamesQueueKey)
	q := store.NewQueue(k.GetCodec(), qStore)
	err := k.checkGames(ctx, q)
	assert.Nil(t, err)
}

func Test_returnFunds(t *testing.T) {
	ctx, votes, k := fakeValidationGame()

	// get the gameID
	qStore := ctx.KVStore(k.activeGamesQueueKey)
	q := store.NewQueue(k.GetCodec(), qStore)
	var gameID int64
	q.Peek(&gameID)

	vote := votes.falseVotes[1]

	initialBalance := k.bankKeeper.GetCoins(ctx, vote.Creator())
	assert.Equal(t, "1000trudex", initialBalance.String())

	err := k.returnFunds(ctx, gameID)
	assert.Nil(t, err)

	// no, not the sneaker
	// FYI - Steve Jobs used to wear NB
	expectedNewBalance := sdk.Coins{vote.Amount()}.Plus(initialBalance)
	actualNewBalance := k.bankKeeper.GetCoins(ctx, vote.Creator())

	assert.Equal(t, expectedNewBalance, actualNewBalance)
}
