package category

import (
	"testing"

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

	catID := ck.Create(ctx, "dog memes", "doggo", "category for dog memes")
	cat, _ := ck.GetCategory(ctx, catID)

	assert.Equal(t, cat.Denom(), "doggo", "should return coin name")
}
