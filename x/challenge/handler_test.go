package challenge

import (
	"encoding/binary"
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestStartChallengeMsg(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	msg := NewStartChallengeMsg(storyID, amount, argument, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	// spew.Dump(res)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(1), x, "incorrect result data")
}

func TestStartChallengeMsg_ErrInsufficientFunds(t *testing.T) {
	ctx, k, sk, ck, _ := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	msg := NewStartChallengeMsg(storyID, amount, argument, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(0), x, "incorrect result data")
}

func TestStartChallengeMsg_ErrInsufficientChallengeAmount(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(1))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	msg := NewStartChallengeMsg(storyID, amount, argument, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(0), x, "incorrect result data")
}

func TestUpdateChallengeMsg(t *testing.T) {
	ctx, k, sk, ck, bankKeeper := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(55))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 3})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give users some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	msg := NewStartChallengeMsg(storyID, amount, argument, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)

	updateMsg := NewUpdateChallengeMsg(1, amount, creator2)

	res = h(ctx, updateMsg)
	spew.Dump(res)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(1), x, "incorrect result data")
}
