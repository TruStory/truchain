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

	title := "Fungi"
	creator := sdk.AccAddress([]byte{1, 2})
	slug := "fungi"
	desc := ""

	msg := NewCreateCategoryMsg(title, creator, slug, desc)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	err := json.Unmarshal(res.Data, &idres)
	assert.NoError(t, err)

	res1 := h(ctx, msg)
	idres1 := new(types.IDResult)
	err = json.Unmarshal(res1.Data, &idres1)
	assert.NoError(t, err)

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
