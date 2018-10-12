package category

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategoryMsg(t *testing.T) {
	validTitle := "Flying cars"
	validCreator := sdk.AccAddress([]byte{1, 2})
	validSlug := "flying-cars"
	validDescription := ""

	msg := NewCreateCategoryMsg(validTitle, validCreator, validSlug, validDescription)
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, "category", msg.Type(), "should get type")
	assert.Equal(t, "create_category", msg.Name(), "should get name")
}

func TestCreateCategoryMsg_Invalid(t *testing.T) {
	invalidTitle := "Flying cars and a bunch of other stuff"
	validCreator := sdk.AccAddress([]byte{1, 2})
	validSlug := "flying-cars"
	validDescription := ""

	msg := NewCreateCategoryMsg(invalidTitle, validCreator, validSlug, validDescription)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCategoryMsg("").Code(), err.Code(), "should throw error")
}
