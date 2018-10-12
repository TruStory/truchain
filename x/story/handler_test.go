package story

import (
	"encoding/binary"
	"strings"
	"testing"

	c "github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestSubmitStoryMsg(t *testing.T) {
	ctx, sk, ck := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	cat := createFakeCategory(ctx, ck)

	body := "fake story"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	msg := NewSubmitStoryMsg(body, cat.ID, creator, kind)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	res1 := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	x1, _ := binary.Varint(res1.Data)
	assert.Equal(t, int64(1), x, "incorrect result data")
	assert.Equal(t, int64(2), x1, "incorrect result data")
}

func TestSubmitStoryMsg_ErrInvalidCategory(t *testing.T) {
	ctx, sk, _ := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	catID := int64(5)

	body := "fake story"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	msg := NewSubmitStoryMsg(body, catID, creator, kind)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInvalidCategory := strings.Contains(res.Log, "801")
	assert.True(t, hasInvalidCategory, "should return err code")
}

func TestByzantineMsg(t *testing.T) {
	ctx, sk, _ := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	fakeMsg := c.NewCreateCategoryMsg("title", sdk.AccAddress([]byte{1, 2}), "slug", "")

	res := h(ctx, fakeMsg)
	hasUnrecognizedMessage := strings.Contains(res.Log, "65542")
	assert.True(t, hasUnrecognizedMessage, "should return err code")
}
