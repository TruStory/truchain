package story

import (
	"encoding/json"
	"strings"
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

	body := "fake story"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	msg := NewSubmitStoryMsg(body, cat.ID, creator, kind)
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
