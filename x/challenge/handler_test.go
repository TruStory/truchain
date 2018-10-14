package challenge

import (
	"encoding/binary"
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestStartChallengeMsg(t *testing.T) {
	ctx, k, sk, ck := mockDB()

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
	assert.Equal(t, int64(1), x, "incorrect result data")
}
