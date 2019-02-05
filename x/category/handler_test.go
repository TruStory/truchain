package category

import (
	"encoding/json"
	"testing"

	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateCategoryMsg(t *testing.T) {
	ctx, ck := mockDB()
	h := NewHandler(ck)
	assert.NotNil(t, h)

	title := "Flying cars"
	creator := sdk.AccAddress([]byte{1, 2})
	slug := "flying-cars"
	desc := ""

	msg := NewCreateCategoryMsg(title, creator, slug, desc)
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

func TestByzantineMsg(t *testing.T) {
	ctx, ck := mockDB()

	h := NewHandler(ck)
	assert.NotNil(t, h)

	res := h(ctx, nil)

	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)

}
