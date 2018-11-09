package challenge

import (
	"encoding/binary"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestSubmitChallengeMsg(t *testing.T) {
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

	msg := NewSubmitChallengeMsg(storyID, amount, argument, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID, "incorrect result data")
}

func TestSubmitChallengeMsg_ErrInsufficientFunds(t *testing.T) {
	ctx, k, sk, ck, _ := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("testcoin", sdk.NewInt(5))
	argument := "test argument"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	msg := NewSubmitChallengeMsg(storyID, amount, argument, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(0), x, "incorrect result data")
}

func TestSubmitChallengeMsg_ErrInsufficientChallengeAmount(t *testing.T) {
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

	msg := NewSubmitChallengeMsg(storyID, amount, argument, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(0), x, "incorrect result data")
}
