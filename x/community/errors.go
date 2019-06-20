package community

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Community errors reserve 800 ~ 899.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeCommunityNotFound   sdk.CodeType = 801
	ErrorCodeInvalidCommunityMsg sdk.CodeType = 802
	ErrorCodeJSONParsing         sdk.CodeType = 803
)

// ErrCommunityNotFound throws an error when the searched category is not found
func ErrCommunityNotFound(id uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeCommunityNotFound, fmt.Sprintf("Community not found with ID: %d", id))
}

// ErrInvalidCommunityMsg throws an error when the searched category is not found
func ErrInvalidCommunityMsg(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidCommunityMsg, fmt.Sprintf("Invalid community msg. Reason: %s", message))
}

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeJSONParsing,
		"JSON parsing error: "+err.Error())
}
