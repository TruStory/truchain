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

	body := "fake story body with minimum length"
	creator := sdk.AccAddress([]byte{1, 2})
	kind := Default
	source := "http://trustory.io"
	evidence := []string{"http://shanesbrain.net"}
	argument := "argument body"

	msg := NewSubmitStoryMsg(argument, body, cat.ID, creator, evidence, source, kind)
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
	evidence := []string{}
	argument := ""

	msg := NewSubmitStoryMsg(argument, body, cat.ID, creator, evidence, source, kind)
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
	evidence := []string{"http://shanesbrain.net"}
	argument := "argument body"
	msg := NewSubmitStoryMsg(argument, body, catID, creator, evidence, source, kind)
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

func TestAddArgumentMsg(t *testing.T) {
	ctx, sk, ck := mockDB()

	h := NewHandler(sk)
	creator := sdk.AccAddress([]byte{1, 2})
	argument := "this is an argument"

	storyID := createFakeStory(ctx, sk, ck)

	msg := NewAddArgumentMsg(storyID, creator, argument)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	story, _ := sk.Story(ctx, storyID)
	assert.Equal(t, story.Arguments[0].Body, argument)

	assert.Equal(t, int64(1), idres.ID, "incorrect result data")
}

func TestAddEvidenceMsg(t *testing.T) {
	ctx, sk, ck := mockDB()

	h := NewHandler(sk)
	creator := sdk.AccAddress([]byte{1, 2})
	evidence := "http://shanesbrain.net"

	storyID := createFakeStory(ctx, sk, ck)

	msg := NewAddEvidenceMsg(storyID, creator, evidence)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID, "incorrect result data")
}

func TestFlagStoryMsg(t *testing.T) {
	ctx, sk, ck := mockDB()

	h := NewHandler(sk)
	creator := sdk.AccAddress([]byte{1, 2})

	storyID := createFakeStory(ctx, sk, ck)

	msg := NewFlagStoryMsg(storyID, creator)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID, "incorrect result data")

	story, _ := sk.Story(ctx, storyID)
	assert.True(t, story.Flagged)
}
