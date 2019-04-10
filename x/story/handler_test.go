package story

import (
	"encoding/json"
	"testing"

	"github.com/TruStory/truchain/types"
	c "github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestSubmitStoryMsg(t *testing.T) {
	ctx, sk, ck := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	cat := createFakeCategory(ctx, ck)

	body := "fake story body with minimum length"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	source := "http://trustory.io"

	msg := NewSubmitStoryMsg(body, cat.ID, creator, source, kind)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	res1 := h(ctx, msg)
	idres := new(types.IDResult)
	idres1 := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)
	_ = json.Unmarshal(res1.Data, &idres1)

	assert.Equal(t, int64(1), idres.ID, "incorrect result data")
	assert.Equal(t, int64(2), idres1.ID, "incorrect result data")
}

func TestSubmitStoryWithoutHostInSourceURLMsg(t *testing.T) {
	ctx, sk, ck := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	cat := createFakeCategory(ctx, ck)

	body := "fake story body with minimum length"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	source := "www.nbd.com"

	msg := NewSubmitStoryMsg(body, cat.ID, creator, source, kind)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	res1 := h(ctx, msg)
	idres := new(types.IDResult)
	idres1 := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)
	_ = json.Unmarshal(res1.Data, &idres1)

	assert.Equal(t, int64(1), idres.ID, "incorrect result data")
	assert.Equal(t, int64(2), idres1.ID, "incorrect result data")
}

func TestSubmitStoryMsgWithOnlyRequiredFields(t *testing.T) {
	ctx, sk, ck := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	cat := createFakeCategory(ctx, ck)

	body := "fake story body with minimum length"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	source := "http://trustory.io"

	msg := NewSubmitStoryMsg(body, cat.ID, creator, source, kind)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	res1 := h(ctx, msg)
	idres := new(types.IDResult)
	idres1 := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)
	_ = json.Unmarshal(res1.Data, &idres1)

	assert.Equal(t, int64(1), idres.ID, "incorrect result data")
	assert.Equal(t, int64(2), idres1.ID, "incorrect result data")
}

func TestSubmitStoryMsg_ErrInvalidCategory(t *testing.T) {
	ctx, sk, _ := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	catID := int64(5)

	body := "fake story body with minimum length"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	source := "http://trustory.io"
	msg := NewSubmitStoryMsg(body, catID, creator, source, kind)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	assert.Equal(t, c.CodeInvalidCategory, res.Code)
	assert.Equal(t, c.DefaultCodespace, res.Codespace)
}

func TestByzantineMsg(t *testing.T) {
	ctx, sk, _ := mockDB()

	h := NewHandler(sk)
	assert.NotNil(t, h)

	fakeMsg := c.NewCreateCategoryMsg("title", sdk.AccAddress([]byte{1, 2}), "slug", "")

	res := h(ctx, fakeMsg)
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)
}
