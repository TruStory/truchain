package community

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Community errors reserve 800 ~ 899.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeCommunityNotFound    sdk.CodeType = 801
	ErrorCodeInvalidCommunityMsg  sdk.CodeType = 802
	ErrorCodeAddressNotAuthorised sdk.CodeType = 803
	ErrorCodeJSONParsing          sdk.CodeType = 804
)

// ErrCommunityNotFound throws an error when the searched category is not found
func ErrCommunityNotFound(id string) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeCommunityNotFound, fmt.Sprintf("Community not found with id: %s", id))
}

// ErrInvalidCommunityMsg throws an error when the searched category is not found
func ErrInvalidCommunityMsg(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeInvalidCommunityMsg, fmt.Sprintf("Invalid community msg. Reason: %s", message))
}

// ErrAddressNotAuthorised throws an error when the address is not admin
func ErrAddressNotAuthorised() sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeAddressNotAuthorised, "This address is not authorised to perform this action.")
}

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace, ErrorCodeJSONParsing, "JSON parsing error: "+err.Error())
}
