package truchain

import (
	"encoding/binary"
	"testing"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TruStory/truchain/x/truchain/db"
	"github.com/stretchr/testify/assert"
)

func TestSubmitStoryMsg(t *testing.T) {
	ctx, _, _, k := db.MockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	body := "fake story"
	cat := ts.DEX
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := ts.Default
	msg := ts.NewSubmitStoryMsg(body, cat, creator, storyType)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(1), x, "incorrect result data")
}
