package claim

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CodeType is the error code type for the module
type CodeType = sdk.CodeType

// Claim error types
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidBody                 CodeType = 101
	ErrorCodeInvalidID                   CodeType = 102
	ErrorCodeNotFound                    CodeType = 103
	ErrorCodeInvalidSType                CodeType = 106
	ErrorCodeClaimsWithCommunityNotFound CodeType = 107
	ErrorCodeInvalidSourceURL            CodeType = 108
	ErrorCodeCreatorJailed               CodeType = 109
	ErrorCodeJSONParsing                 CodeType = 110
)

// ErrInvalidBody throws an error on invalid claim body
func ErrInvalidBody(body string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBody,
		"Invalid claim body: "+body)
}

// ErrInvalidID throws an error on invalid claim body
func ErrInvalidID() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidID,
		"Invalid claim ID")
}

// ErrInvalidCommunityID throws an error on invalid community id
func ErrInvalidCommunityID(id uint64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBody,
		fmt.Sprintf("Invalid community id: %d", id))
}

// ErrInvalidSourceURL throws an error when a URL in invalid
func ErrInvalidSourceURL(url string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidSourceURL,
		"Invalid source URL: "+url)
}

// ErrCreatorJailed throws an error on jailed creator
func ErrCreatorJailed(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeCreatorJailed,
		"Creator cannot be jailed: "+addr.String())
}

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeJSONParsing,
		"JSON parsing error: "+err.Error())
}
