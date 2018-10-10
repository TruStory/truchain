package category

import (
	"encoding/binary"
	"strings"
	"testing"

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
	x, _ := binary.Varint(res.Data)
	x1, _ := binary.Varint(res1.Data)
	assert.Equal(t, int64(1), x, "incorrect result data")
	assert.Equal(t, int64(2), x1, "incorrect result data")
}

func TestByzantineMsg(t *testing.T) {
	ctx, ck := mockDB()

	h := NewHandler(ck)
	assert.NotNil(t, h)

	res := h(ctx, nil)
	hasUnrecognizedMessage := strings.Contains(res.Log, "65542")
	assert.True(t, hasUnrecognizedMessage, "should return err code")
}
