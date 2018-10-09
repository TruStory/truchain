package truchain

import (
	"encoding/binary"
	"strings"
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TruStory/truchain/x/truchain/db"
	"github.com/stretchr/testify/assert"
)

func TestSubmitStoryMsg(t *testing.T) {
	ctx, _, _, k := db.MockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	cat := db.CreateFakeCategory(ctx, k)

	body := "fake story"
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := ts.Default
	msg := ts.NewSubmitStoryMsg(body, cat.ID, creator, storyType)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	res1 := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	x1, _ := binary.Varint(res1.Data)
	assert.Equal(t, int64(1), x, "incorrect result data")
	assert.Equal(t, int64(2), x1, "incorrect result data")
}

func TestSubmitStoryMsg_ErrInvalidCategory(t *testing.T) {
	ctx, _, _, k := db.MockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	catID := int64(5)

	body := "fake story"
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := ts.Default
	msg := ts.NewSubmitStoryMsg(body, catID, creator, storyType)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInvalidCategory := strings.Contains(res.Log, "708")
	assert.True(t, hasInvalidCategory, "should return err code")
}

func TestCreateCategoryMsg(t *testing.T) {
	ctx, _, _, k := db.MockDB()
	h := NewHandler(k)
	assert.NotNil(t, h)

	title := "Flying cars"
	creator := sdk.AccAddress([]byte{1, 2})
	slug := "flying-cars"
	desc := ""

	msg := ts.NewCreateCategoryMsg(title, creator, slug, desc)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	res1 := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	x1, _ := binary.Varint(res1.Data)
	assert.Equal(t, int64(1), x, "incorrect result data")
	assert.Equal(t, int64(2), x1, "incorrect result data")
}

func TestBackStoryMsg_FailBasicValidation(t *testing.T) {
	ctx, _, _, k := db.MockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount, _ := sdk.ParseCoin("5trushane")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 5 * time.Hour
	msg := ts.NewBackStoryMsg(storyID, amount, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInvalidBackingPeriod := strings.Contains(res.Log, "706")
	assert.True(t, hasInvalidBackingPeriod, "should return err code")
}

func TestBackStoryMsg_FailInsufficientFunds(t *testing.T) {
	ctx, _, _, k := db.MockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount, _ := sdk.ParseCoin("5trushane")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	msg := ts.NewBackStoryMsg(storyID, amount, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInsufficientFunds := strings.Contains(res.Log, "65541")
	assert.True(t, hasInsufficientFunds, "should return err code")
}

func TestBackStoryMsg(t *testing.T) {
	ctx, ms, am, k := db.MockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := db.CreateFakeStory(ms, k)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := db.CreateFakeFundedAccount(ctx, am, sdk.Coins{amount})
	duration := 99 * time.Hour
	msg := ts.NewBackStoryMsg(storyID, amount, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	x, _ := binary.Varint(res.Data)
	assert.Equal(t, int64(1), x, "incorrect result backing id")
}

func TestByzantineMsg(t *testing.T) {
	ctx, _, _, k := db.MockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	fakeMsg := ts.NewAddCommentMsg(int64(5), "test", sdk.AccAddress([]byte{1, 2}))

	res := h(ctx, fakeMsg)
	hasUnrecognizedMessage := strings.Contains(res.Log, "65542")
	assert.True(t, hasUnrecognizedMessage, "should return err code")
}
