package db

import (
	"testing"

	ts "github.com/TruStory/truchain/x/truchain/types"

	"github.com/stretchr/testify/assert"
)

func TestGetCategory_ErrBackingNotFound(t *testing.T) {
	ctx, _, _, k := MockDB()
	id := int64(5)

	_, err := k.GetCategory(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ts.ErrCategoryNotFound(id).Code(), err.Code(), "should get error")
}

func TestGetCategory(t *testing.T) {
	ctx, _, _, k := MockDB()

	catID, _ := k.NewCategory(ctx, "dog memes", "doggo", "category for dog memes")
	cat, _ := k.GetCategory(ctx, catID)

	assert.Equal(t, cat.CoinName(), "doggo", "should return coin name")
}
