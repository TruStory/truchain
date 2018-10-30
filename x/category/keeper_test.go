package category

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCategory_ErrCategoryNotFound(t *testing.T) {
	ctx, ck := mockDB()
	id := int64(5)

	_, err := ck.GetCategory(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrCategoryNotFound(id).Code(), err.Code(), "should get error")
}

func TestNewGetCategory(t *testing.T) {
	ctx, ck := mockDB()

	catID, _ := ck.NewCategory(ctx, "dog memes", sdk.AccAddress([]byte{1, 2}), "doggo", "category for dog memes")
	cat, _ := ck.GetCategory(ctx, catID)

	assert.Equal(t, cat.CoinName(), "doggo", "should return coin name")
}

func TestInitCategories(t *testing.T) {
	ctx, k := mockDB()

	categories := map[string]string{
		"btc":      "Bitcoin",
		"shitcoin": "Shitcoins",
	}

	creator := sdk.AccAddress([]byte{1, 2})

	err := k.InitCategories(ctx, creator, categories)
	assert.Nil(t, err)

	cat, _ := k.GetCategory(ctx, 1)
	assert.Contains(t, categories, cat.Slug)

	cat, _ = k.GetCategory(ctx, 2)
	assert.Contains(t, categories, cat.Slug)
}
