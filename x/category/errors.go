package category

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Category errors reserve 800 ~ 899.
const (
	DefaultCodespace sdk.CodespaceType = "category"

	CodeInvalidCategory    sdk.CodeType = 801
	CodeCategoryNotFound   sdk.CodeType = 802
	CodeInvalidCategoryMsg sdk.CodeType = 803
	CodeDenomMismatch      sdk.CodeType = 804
)

// ErrInvalidCategory throws an error when the category is invalid
func ErrInvalidCategory(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidCategory, "Invalid story category: "+fmt.Sprintf("%d", id))
}

// ErrInvalidCategoryMsg throws an error when a category msg is invalid
func ErrInvalidCategoryMsg(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidCategoryMsg, msg)
}

// ErrCategoryNotFound throws an error when the searched category is not found
func ErrCategoryNotFound(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCategoryNotFound, "Category id not found: "+fmt.Sprintf("%d", id))
}

// ErrCodeDenomMismatch throws an error when a category msg is invalid
func ErrCodeDenomMismatch(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDenomMismatch, "Coins could not be added to category because the coin denom does not match the category denom.")
}
