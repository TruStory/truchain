package challenge

import (
	"net/url"
	"testing"

	store "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, k, _, _, _ := mockDB()

	tags := k.NewResponseEndBlock(ctx)
	assert.Equal(t, sdk.Tags{}, tags)
}

func Test_checkExpiredGames(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	_, err := k.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	q := store.NewQueue(k.GetCodec(), k.GetStore(ctx))
	err = checkExpiredGames(ctx, k, q)
	assert.Nil(t, err)
}

func Test_returnFunds(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	id, err := k.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	game, _ := k.Get(ctx, id)

	err = returnFunds(ctx, k, game.ID)
	assert.Nil(t, err)

	coin := bankKeeper.GetCoins(ctx, creator)
	assert.Equal(t, sdk.Coins{amount}, coin)
}
